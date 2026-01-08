package consumer

import (
	"github.com/IBM/sarama"
)

type ConsumerGroupHandler struct {
	Messages chan<- *sarama.ConsumerMessage
}

func (h *ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.Messages <- msg
		session.MarkMessage(msg, "")
	}
	return nil
}
