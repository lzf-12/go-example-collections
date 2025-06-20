package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/lzf-12/go-example-collections/internal/config"
	handler "github.com/lzf-12/go-example-collections/internal/consumer/handler"
	pubsub "github.com/lzf-12/go-example-collections/internal/consumer/model"
	"github.com/lzf-12/go-example-collections/msgbroker/adapter/kafka"
)

func InitKafkaConsumer(ctx context.Context) error {

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Printf("load config failed: %v", err)
		return err
	}

	kafkaBrokerServer := fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)
	consumerGroupId := cfg.KafkaConsumerGroupID

	// admin configuration map
	admincfg := kafka.NewKafkaConfigMap()
	admincfg.Set(fmt.Sprintf("bootstrap.servers=%s", kafkaBrokerServer))

	// consumer configuration map
	consumercfg := kafka.NewKafkaConfigMap()
	consumercfg.Set(fmt.Sprintf("bootstrap.servers=%s", kafkaBrokerServer))
	consumercfg.Set(fmt.Sprintf("group.id=%s", consumerGroupId))
	consumercfg.Set(fmt.Sprintf("auto.offset.reset=%s", "earliest"))

	log.Println("initialize kafka consumer client")
	kc, err := kafka.NewKafkaConsumerClient(consumercfg, admincfg)
	if err != nil {
		log.Println("error initialize kafka consumer client: ", err)
	}

	// healtcheck
	log.Println("kafka healthcheck...")
	err = kc.HealthCheck(ctx)
	if err != nil {
		log.Println("healtcheck kafka error: ", err)
		return err
	}
	log.Println("kafka ok")

	// TODO: centralize topic partition configuration in .yml
	// topic handlers map
	topicHandlers := []kafka.TopicHandler{
		{Topic: pubsub.TopicOrderV2Json, Handler: handler.OrderHandlerV2Json, Partitions: 1, ReplicationFactor: 1},
		{Topic: pubsub.TopicOrderV2Xml, Handler: handler.OrderHandlerV2Xml, Partitions: 1, ReplicationFactor: 1},
	}

	// subscribe all topic and handlers
	err = kc.SubscribeTopics(ctx, topicHandlers)
	if err != nil {
		log.Println("subscribe single topic error: ", err)
		return err
	}

	log.Println("success initialize kafka consumer client")
	return nil
}
