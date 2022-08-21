package kafka

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/go-zookeeper/zk"
)

func NewClient(add ...string) (sarama.Client, error) {
	config := sarama.NewConfig()
	return sarama.NewClient(add, config)
}

type Brokder struct {
	Id   string
	Host string
	Port int
}

// Addr return the host:port as a address
func (b Brokder) Addr() string {
	return fmt.Sprintf("%s:%d", b.Host, b.Port)
}

const (
	ZKAdd       = "localhost:2181"
	BrokderPath = "/brokers/ids"
)

func GetBrokders(zkAddr string) ([]Brokder, error) {
	conn, _, err := zk.Connect([]string{zkAddr}, 5*time.Second)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	ids, _, err := conn.Children(BrokderPath)
	if err != nil {
		return nil, err
	}
	brokders := make([]Brokder, 0, len(ids))

	for _, id := range ids {
		brokderInfo, _, err := conn.Get(fmt.Sprintf("%s/%s", BrokderPath, id))
		if err != nil {
			continue
		}
		broker := Brokder{
			Id: id,
		}
		err = json.Unmarshal(brokderInfo, &broker)
		if err != nil {
			continue
		}
		brokders = append(brokders, broker)
	}
	return brokders, nil
}
