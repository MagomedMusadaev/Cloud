package kafka

import (
	"Cloud/logger"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"net/http"
)

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

	// Создание нового продюсера с указанием адреса Kafka
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		//"acks":               "all",         // Подтверждение после полной репликации
		//"retries":            3,             // Повторные попытки в случае неудачи
		//"enable.idempotence": true,          // Включение идемпотентности для предотвращения дублирования
	})
	if err != nil {
		// Если произошла ошибка при создании продюсера, логируем ошибку и выходим из функции
		logger.Error("Error creating producer:" + err.Error())
		return
	}
	defer producer.Close()

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	// Отправка сообщения в указанный топик
	err = producer.Produce(&kafka.Message{
		// Указание топика и партиции, в которую будет отправлено сообщение
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		logger.Error("Error producing message:" + err.Error())
		return
	}

	// Ожидание события доставки и получение результата
	e := <-deliveryChan
	msg := e.(*kafka.Message) // Приведение типа события к сообщению Kafka

	// Проверка на наличие ошибок при доставке сообщения
	if msg.TopicPartition.Error != nil {
		logger.Error("Delivery failed:" + msg.TopicPartition.Error.Error())
	} else {
		logger.Info(fmt.Sprintf("Message delivered to %v\n", msg.TopicPartition))
	}
}
