package kafka

import (
	"fmt"
)

func Info() {

	client, err := NewClientFromZK(ZKAdd)
	if err != nil {
		logger.Fatal(err)
	}
	brokers := client.Brokers()
	defer func() {
		for _, broker := range brokers {
			broker.Close()
		}
	}()
	brokersName := make([]string, 0, len(brokers))
	for _, broker := range brokers {
		brokersName = append(brokersName, fmt.Sprintf("%d-%s", broker.ID(), broker.Addr()))
	}
	topics, err := client.Topics()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Printf("client has brokers:%v, has topics:%v\n", brokersName, topics)
}
