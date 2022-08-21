package kafka

import (
	"fmt"
)

func Info() {
	brokerInfos, err := GetBrokders(ZKAdd)
	if err != nil {
		logger.Fatal(err)
	}
	if len(brokerInfos) == 0 {
		logger.Error("cannot get any brokder")
	}
	adds := make([]string, 0, len(brokerInfos))
	for _, b := range brokerInfos {
		adds = append(adds, b.Addr())
	}

	client, err := NewClient(adds...)
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
	logger.Printf("client(%s) has brokers:%v, has topics:%v\n", brokerInfos[0].Addr(), brokersName, topics)
}
