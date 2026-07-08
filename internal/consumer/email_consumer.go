package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/ESE-MONDAY/relay-service/internal/event"
	"github.com/ESE-MONDAY/relay-service/internal/processor"
)

type EmailConsumer struct {
	reader    *kafka.Reader
	processor *processor.EmailProcessor
}

func NewEmailConsumer(
	brokers []string,
	topic string,
	group string,
	p *processor.EmailProcessor,
) *EmailConsumer {

	return &EmailConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: group,
		}),
		processor: p,
	}
}

func (c *EmailConsumer) Start(ctx context.Context) {
	log.Println("Kafka consumer started")

	defer func() {
		log.Println("Closing Kafka consumer...")
		_ = c.reader.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer shutting down...")
			return

		default:
			msg, err := c.reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}

				log.Println("kafka read error:", err)
				continue
			}

			var ev event.EmailEvent

			if err := json.Unmarshal(msg.Value, &ev); err != nil {
				log.Println("invalid event:", err)
				continue
			}

			switch ev.Type {

			case event.EmailEventTypeCreated,
				event.EmailEventTypeRetry:

				if err := c.processor.Process(
					ctx,
					ev.EmailID,
					ev.Retry,
				); err != nil {
					log.Println("email processing failed:", err)
				}

			case event.EmailEventTypeFailed:
				log.Printf(
					"received failed event for email %s\n",
					ev.EmailID,
				)

			case event.EmailEventTypeDLQ:
				log.Printf(
					"received dead-letter event for email %s\n",
					ev.EmailID,
				)

			default:
				log.Printf(
					"unknown event type: %s\n",
					ev.Type,
				)
			}
		}
	}
}
