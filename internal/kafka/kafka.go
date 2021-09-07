package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/ozonva/ova-purchase-api/internal/config"
	"github.com/rs/zerolog/log"
)

type Action int

const (
	CreatePurchase Action = iota
	MultiCreatePurchase
	UpdatePurchase
	RemovePurchase
)

type message struct {
	Id     uuid.UUID   `json:"id"`
	Action Action      `json:"action"`
	Value  interface{} `json:"value"`
}

type KafkaProducer interface {
	Send(message message) error
}

type producer struct {
	saramaProducer sarama.SyncProducer
	topic          string
}

func NewMessage(action Action, value interface{}) message {
	id, _ := uuid.NewRandom()
	return message{
		Id:     id,
		Action: action,
		Value:  value,
	}
}

func NewProducer(config config.KafkaConfiguration) (KafkaProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	saramaConfig.Producer.Return.Successes = true

	saramaProducer, err := sarama.NewSyncProducer(config.Brokers, saramaConfig)

	if err != nil {
		log.Fatal().Err(err).Msg("Kafka producer: failed to create")
		return nil, err
	}

	return &producer{
		saramaProducer: saramaProducer,
		topic:          config.Topic,
	}, nil
}

func (s *producer) Send(message message) error {
	jsonMes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, _, err = s.saramaProducer.SendMessage(
		&sarama.ProducerMessage{
			Topic: s.topic,
			Key:   sarama.StringEncoder(message.Id.String()),
			Value: sarama.StringEncoder(jsonMes),
		})
	return err
}

func (s *producer) Disposal() {
	log.Debug().Msg("Close kafka producer...")
	if err := s.saramaProducer.Close(); err != nil {
		log.Error().Msgf("Failed close kafka producer %v", err)
	}
}
