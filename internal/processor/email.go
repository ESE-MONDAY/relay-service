package processor

import (
	"context"

	"github.com/ESE-MONDAY/relay-service/internal/models"
	"github.com/ESE-MONDAY/relay-service/internal/queue"
	"github.com/ESE-MONDAY/relay-service/internal/repository"
	"github.com/ESE-MONDAY/relay-service/internal/sender"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type EmailProcessor struct {
	repo   repository.EmailStore
	sender sender.Sender
	log    *zap.Logger
}

func NewEmailProcessor(
	repo repository.EmailStore,
	sender sender.Sender,
	log *zap.Logger,
) *EmailProcessor {

	return &EmailProcessor{
		repo:   repo,
		sender: sender,
		log:    log,
	}
}

func (p *EmailProcessor) Process(
	ctx context.Context,
	job queue.Job,
) error {

	email, err := p.loadEmail(ctx, job.EmailID)
	if err != nil {
		return err
	}

	if err := p.markProcessing(ctx, email.ID); err != nil {
		return err
	}

	p.log.Info("calling sender")

	if err := p.sender.Send(ctx, email); err != nil {
		_ = p.markFailed(ctx, email.ID)
		return err
	}

	p.log.Info("sender returned")

	if err := p.markSent(ctx, email.ID); err != nil {
		return err
	}

	p.log.Info("processing complete")

	return nil
}
func (p *EmailProcessor) loadEmail(
	ctx context.Context,
	id uuid.UUID,
) (*models.Email, error) {

	email, err := p.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	p.log.Info(
		"email loaded",
		zap.String("email_id", email.ID.String()),
		zap.String("from", email.From),
		zap.String("to", email.To),
		zap.String("subject", email.Subject),
	)

	return email, nil
}
func (p *EmailProcessor) markProcessing(
	ctx context.Context,
	id uuid.UUID,
) error {

	p.log.Info(
		"marking email as processing",
		zap.String("email_id", id.String()),
	)

	return p.repo.UpdateStatus(
		ctx,
		id,
		models.EmailProcessing,
	)
}

func (p *EmailProcessor) markSent(
	ctx context.Context,
	id uuid.UUID,
) error {

	p.log.Info(
		"email sent",
		zap.String("email_id", id.String()),
	)

	return p.repo.UpdateStatus(
		ctx,
		id,
		models.EmailSent,
	)
}
func (p *EmailProcessor) markFailed(
	ctx context.Context,
	id uuid.UUID,
) error {

	p.log.Warn(
		"email failed",
		zap.String("email_id", id.String()),
	)

	return p.repo.UpdateStatus(
		ctx,
		id,
		models.EmailFailed,
	)
}
