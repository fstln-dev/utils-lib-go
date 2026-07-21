package rabbit

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type producter struct {
	conn     *amqpConn
	confir   chan amqp.Confirmation
	exchange string
}

type Producter interface {
	Headers(h *Header) ProducerOption
	MessageID(id string) ProducerOption
	Publish(ctx context.Context, routekey string, body []byte, opts ...ProducerOption) error
	Confirm(ctx context.Context, ack, nack func(uint64)) error
}

//
func (a *amqpConn) NewProducer(exchange string) (Producter, error) {

	// 检测
	channel, err := a.getChannel()
	if err != nil {
		return nil, err
	}

	err = channel.ExchangeDeclarePassive(
		exchange, "topic",
		false, false, false, false,
		amqp.Table{})
	if err != nil {
		return nil, err
	}

	return &producter{
		conn:     a,
		exchange: exchange,
	}, nil
}

type ProducerOption func(p *amqp.Publishing)

func (p *producter) Headers(h *Header) ProducerOption {
	return func(p *amqp.Publishing) {
		p.Headers = h.t
	}
}

func (p *producter) MessageID(id string) ProducerOption {
	return func(p *amqp.Publishing) {
		p.MessageId = id
	}
}

func (p *producter) Publish(ctx context.Context, routekey string, body []byte, opts ...ProducerOption) error {

	channel, err := p.conn.getChannel()
	if err != nil {
		return err
	}

	publishing := amqp.Publishing{
		Headers:         amqp.Table{},
		ContentType:     JSON,
		ContentEncoding: UTF8,
		Body:            body,
		DeliveryMode:    amqp.Persistent, // 消息持久化模式
		Priority:        0,
	}

	for _, o := range opts {
		o(&publishing)
	}

	return channel.PublishWithContext(ctx, p.exchange, routekey, true, false, publishing)
}

func (p *producter) Confirm(ctx context.Context, ack, nack func(uint64)) error {
	c, err := p.conn.getChannel()
	if err != nil {
		return err
	}

	// 必须声明为 Confirm(false) 才会触发 ack 的 notify
	if ack != nil {
		if err := p.conn.channel.Confirm(false); err != nil {
			return err
		}
	}

	if p.confir != nil {
		p.confir = make(chan amqp.Confirmation)
	}

	ff := c.NotifyPublish(p.confir)
	for {
		time.Sleep(100 * time.Millisecond)
		select {
		case conf := <-ff:
			if conf.Ack {
				ack(conf.DeliveryTag)
			} else {
				nack(conf.DeliveryTag)
			}
		case <-ctx.Done():
		}
	}
}
