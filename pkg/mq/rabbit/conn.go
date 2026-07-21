package rabbit

import (
	"context"
	"sync"

	"github.com/go-kratos/kratos/v2/log"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	JSON = "application/json"
	UTF8 = "utf-8"
)

type amqpConn struct {
	amqpChan
	url  string
	conn *amqp.Connection
	log  *log.Helper
}

type amqpChan struct {
	channel *amqp.Channel
	locker  sync.Locker
}

type Header struct {
	t amqp.Table
}

type Delivery amqp.Delivery

type RabbitMQ interface {
	NewProducer(exchange string) (Producter, error)
	Listen(ctx context.Context, queue, consume string, handle func(d Delivery) (requeue bool, err error), opts ...ConsumerOption) (Listener, error)
}

func NewAMQP(url string, logger log.Logger) (RabbitMQ, error) {
	re := &amqpConn{
		url: url,
		amqpChan: amqpChan{
			locker: &sync.Mutex{},
		},
		log: log.NewHelper(logger),
	}

	_, err := re.getChannel()
	if err != nil {
		return nil, err
	}

	return re, nil
}

func (a *amqpConn) getChannel() (*amqp.Channel, error) {

	a.locker.Lock()
	defer a.locker.Unlock()

	var err error

	// conn still alived
	if a.conn == nil || a.conn.IsClosed() {
		// resurrect the conn
		a.conn, err = amqp.Dial(a.url)
		if err != nil {
			a.log.Errorf("resurrect conn err: %v", err)
			return nil, err
		}
	}

	// channel still alived
	if a.channel == nil || a.channel.IsClosed() {
		// resurrect the channel
		a.channel, err = a.conn.Channel()
		if err != nil {
			a.log.Errorf("resurrect channel err: %v", err)
			return nil, err
		}
	}

	// return channel

	return a.channel, nil
}

func NewHeaders(t ...*amqp.Table) *Header {

	re := &Header{
		t: amqp.Table{},
	}

	if t != nil {
		re.t = *t[0]
	}

	return re
}

func (h *Header) Set(key string, value interface{}) error {
	if err := (amqp.Table{key: value}).Validate(); err != nil {
		return err
	}

	h.t[key] = value
	return nil
}

func (h *Header) Get(key string) interface{} {
	return h.t[key]
}
