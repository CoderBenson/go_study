package kafka

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

const Topic = "topic_tedst"

func Product(msg string) {
	client, err := NewClientFromZK(ZKAdd)
	if err != nil {
		logger.Fatal(err)
	}
	defer client.Close()
	producer, err := sarama.NewAsyncProducerFromClient(client)
	if err != nil {
		logger.Fatal(err)
	}
	defer producer.Close()
	producer.Input() <- &sarama.ProducerMessage{
		Topic: Topic,
		Value: sarama.StringEncoder(msg),
	}
}

type Consumer struct {
	consume func(*sarama.ConsumerMessage)
}

func NewConsumer(consume func(*sarama.ConsumerMessage)) *Consumer {
	return &Consumer{consume: consume}
}

func (Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c Consumer) ConsumeClaim(_ sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		c.consume(msg)
	}
	return nil
}

func Consume(ctx context.Context, consume func(*sarama.ConsumerMessage)) {
	client, err := NewClientFromZK(ZKAdd)
	if err != nil {
		logger.Fatal(err)
	}
	defer client.Close()
	if err != nil {
		logger.Fatal(err)
	}
	consumerGroup, err := sarama.NewConsumerGroupFromClient("group", client)
	if err != nil {
		logger.Fatal(err)
	}
	defer consumerGroup.Close()
	run := true
	for run {
		logger.Info("start consume")
		select {
		case <-ctx.Done():
			run = false
		default:
			err := consumerGroup.Consume(ctx, []string{Topic}, NewConsumer(consume))
			if err != nil {
				run = false
			}
		}
		if !run {
			logger.Info("consumer exit")
		}
	}
}

func ProductConsume() {
	baseCtx := context.Background()
	ctx, cancelFunc := context.WithCancel(baseCtx)
	go func() {
		Consume(ctx, func(msg *sarama.ConsumerMessage) {
			logger.Infof("receive msg:%s(%d)\n", string(msg.Value), msg.Timestamp.Unix())
		})
	}()
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	run := true
	go func() {
		<-exit
		logger.Info("exit by signal")
		cancelFunc()
		run = false
	}()
	reader := bufio.NewReader(os.Stdin)
	for run {
		line, _, err := reader.ReadLine()
		if err != nil {
			logger.Fatal(err)
		}
		if run {
			Product(string(line))
		}
		logger.Infof("sned msg:%s\n", line)
	}
	fmt.Println("exit")
}
