package queue

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	Publish(ctx context.Context, u string) error
}

func NewRabbitMqPublisher(queue *amqp.Queue, channel *amqp.Channel) Publisher {
	return &RabbitMqPublisher{queue, channel}
}

type RabbitMqPublisher struct {
	queue   *amqp.Queue
	channel *amqp.Channel
}

func (p *RabbitMqPublisher) Publish(ctx context.Context, u string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := p.channel.PublishWithContext(ctx,
		"",           // exchange
		p.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(u),
		})

	if err != nil {
		return err
	}

	return nil
}
