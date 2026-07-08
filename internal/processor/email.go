package processor

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/ESE-MONDAY/relay-service/internal/event"
	"github.com/ESE-MONDAY/relay-service/internal/models"
	"github.com/ESE-MONDAY/relay-service/internal/queue"
	"github.com/ESE-MONDAY/relay-service/internal/repository"
	"github.com/ESE-MONDAY/relay-service/internal/sender"
)

const maxRetries = 3

type EmailProcessor struct {
	repo      repository.EmailStore
	publisher queue.Publisher
	sender    sender.Sender
	log       *zap.Logger
}

func NewEmailProcessor(
	repo repository.EmailStore,
	publisher queue.Publisher,
	sender sender.Sender,
	log *zap.Logger,
) *EmailProcessor {
	return &EmailProcessor{
		repo:      repo,
		publisher: publisher,
		sender:    sender,
		log:       log,
	}
}

func (p *EmailProcessor) Process(
	ctx context.Context,
	emailID uuid.UUID,
	retry int,
) error {

	// 1. Atomic claim (only one worker wins)
	email, err := p.repo.ClaimForProcessing(ctx, emailID)
	if err != nil {
		return err
	}

	if email == nil {
		p.log.Info("email already processed",
			zap.String("email_id", emailID.String()),
		)
		return nil
	}

	p.log.Info("processing email",
		zap.String("email_id", email.ID.String()),
		zap.Int("retry", retry),
	)

	// 2. Send email
	if err := p.sender.Send(ctx, email); err != nil {

		p.log.Warn("email send failed",
			zap.String("email_id", email.ID.String()),
			zap.Int("retry", retry),
			zap.Error(err),
		)

		// mark failure attempt
		_ = p.repo.UpdateStatus(ctx, email.ID, models.EmailFailed)

		// 3. retry via EVENT (NOT sleep, NOT local loop)
		if retry < maxRetries {

			ev := event.EmailEvent{
				Version: "v1",
				EventID: uuid.NewString(),
				Type:    "email.retry",
				EmailID: email.ID,
				Retry:   retry + 1,
			}

			if err := p.publisher.Publish(ctx, ev); err != nil {
				p.log.Error("failed to publish retry event",
					zap.Error(err),
				)
			}

			return nil
		}

		// 4. dead letter
		_ = p.repo.UpdateStatus(ctx, email.ID, models.EmailDeadLetter)

		p.log.Error("email moved to dead letter",
			zap.String("email_id", email.ID.String()),
		)

		return nil
	}

	// 5. success
	if err := p.repo.UpdateStatus(ctx, email.ID, models.EmailSent); err != nil {
		return err
	}

	p.log.Info("email sent successfully",
		zap.String("email_id", email.ID.String()),
	)

	return nil
}
