package consumer

import (
	"context"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type ConsumerGroupHandler struct {
	WorkerID int
}

// Setup é chamado no início da sessão
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	log.Printf("[Worker-%d] Sessão iniciada", h.WorkerID)
	return nil
}

// Cleanup é chamado no final da sessão
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Printf("[Worker-%d] Sessão finalizada", h.WorkerID)
	return nil
}

// ConsumeClaim é chamado para cada partição atribuída
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("[Worker-%d] Mensagem recebida:", h.WorkerID)
		log.Printf("   Tópico: %s", msg.Topic)
		log.Printf("   Partição: %d", msg.Partition)
		log.Printf("   Offset: %d", msg.Offset)
		log.Printf("   Key: %s", string(msg.Key))
		log.Printf("   Value: %s", string(msg.Value))

		session.MarkMessage(msg, "")
	}
	return nil
}

// RunConsumers cria múltiplos consumers em paralelo dentro do mesmo consumer group
func RunConsumers(ctx context.Context, groupID, topic string, brokers []string, numWorkers int) {
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			config := sarama.NewConfig()
			config.Version = sarama.MaxVersion
			config.Consumer.Return.Errors = true
			config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
			config.Consumer.Offsets.Initial = sarama.OffsetOldest

			consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
			if err != nil {
				log.Fatalf("[Worker-%d] Erro criando consumer group: %v", workerID, err)
			}
			defer consumerGroup.Close()

			// Goroutine para capturar erros
			go func() {
				for err := range consumerGroup.Errors() {
					log.Printf("[Worker-%d] Erro no consumer: %v", workerID, err)
				}
			}()

			handler := &ConsumerGroupHandler{WorkerID: workerID}

			for {
				log.Printf("[Worker-%d] Iniciando consumo no tópico %s...", workerID, topic)
				if err := consumerGroup.Consume(ctx, []string{topic}, handler); err != nil {
					log.Printf("[Worker-%d] Erro consumindo mensagens: %v", workerID, err)
				}

				if ctx.Err() != nil {
					log.Printf("[Worker-%d] Contexto cancelado, encerrando...", workerID)
					return
				}
			}
		}(i)
	}

	wg.Wait()
}
