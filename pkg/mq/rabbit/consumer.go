package rabbit

import (
	"context"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Listener interface {
	Close()
}

type listener struct {
	conn     *amqpConn
	queue    string
	consume  string
	handle   func(d Delivery) (requeue bool, err error)
	stop     context.CancelFunc
	interval time.Duration
}

type ConsumerOption func(l *listener)

func Interval(t time.Duration) ConsumerOption {
	return func(l *listener) {
		l.interval = t
	}
}

func (a *amqpConn) Listen(
	ctx context.Context,
	queue, consume string,
	handle func(d Delivery) (requeue bool, err error),
	opts ...ConsumerOption,
) (Listener, error) {

	l := &listener{
		conn:     a,
		queue:    queue,
		consume:  consume,
		handle:   handle,
		interval: 5 * time.Second,
	}

	for _, o := range opts {
		o(l)
	}

	ctx, l.stop = context.WithCancel(ctx)

	// 异步监听
	go l.createConsumer(ctx)
	// 第一次失败时返回消息（暂无方案）
	return l, nil
}

func (l *listener) createConsumer(ctx context.Context) {
	for {

		// 创建/获取 MQ 通道
		channel, err := l.conn.getChannel()
		if err != nil {

			l.conn.log.Errorf("Listener closed by error: %v", err)

			// sleep 5 second
			time.Sleep(l.interval)

			// try agian
			continue
		}

		// 执行监听
		if err := l.exec(ctx, channel); err != nil {
			l.conn.log.Errorf("Listener exec function cloded by error: %v", err)
			// sleep 5 second
			time.Sleep(l.interval)
			continue
		}

		// 关闭
		l.conn.log.Info("Listener closed")
		return
	}
}

func (l *listener) exec(ctx context.Context, channel *amqp.Channel) error {

	deliverys, err := channel.Consume(
		l.queue, l.consume,
		false, false, false, false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case d, ok := <-deliverys:
			if !ok {
				return errors.New("channel was accidentally closed")
			}

			ifrequeue, err := l.handle(Delivery(d))
			if err != nil {
				// Nack，及错误重推
				// true: 重推, false: 进入死信队列
				for i := 0; i < 4; i++ {
					err = d.Nack(false, ifrequeue)
					if err == nil {
						break
					}
				}
				if err != nil {
					// log
					l.conn.log.Errorf("the message %s Nack Error: %v", d.MessageId, err)
					return err
				}
				continue
			}

			// 消息确认
			for i := 0; i < 4; i++ {
				err = d.Ack(false)
				if err == nil {
					break
				}
			}
			if err != nil {
				// log
				l.conn.log.Errorf("the message %s Ack Error: %v", d.MessageId, err)
				return err
			}
		}
	}
}

func (l *listener) Close() {
	l.stop()
}
