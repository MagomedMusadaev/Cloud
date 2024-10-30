package kafka

import (
	"Cloud/logger"
	"fmt"
	"github.com/IBM/sarama"
	"net/http"
)

// Обработчик для отправки запроса на Kafka
func SendToKafkaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	produceMessage("file_topic", "Работайте братья!")

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Request sent to Kafka"))
}

// Функция для отправки сообщения в указанный топик Kafka
func produceMessage(topic string, message string) {
	// Конфигурация продюсера
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Подтверждение после полной репликации
	config.Producer.Retry.Max = 5                    // Повторные попытки при неудаче
	config.Producer.Idempotent = true                // Включение идемпотентности
	config.Producer.Return.Successes = true          // Ожидание подтверждения успешной отправки
	config.Net.MaxOpenRequests = 1                   // Требуется для идемпотентного продюсера

	// Создание продюсера
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, config)
	if err != nil {
		logger.Error("Error creating producer: " + err.Error())
		return
	}
	defer producer.Close()

	// Создание сообщения
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	// Отправка сообщения и получение результата
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		logger.Error("Error producing message: " + err.Error())
		return
	}

	logger.Info(fmt.Sprintf("Message delivered to partition %d at offset %d\n", partition, offset))
}
