package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

type Comment struct {
	Text string `form:"text" json:"text"`
}

func main() {
	app := gin.Default()
	app.POST("/api/v1/comments", createComment)
	app.Run()
}

func ConnectProducer(brokerUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	conn, err := sarama.NewSyncProducer(brokerUrl, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func PushCommentToQueue(topic string, message []byte) error {
	brokerUrl := []string{"localhost:29092"}
	producer, err := ConnectProducer(brokerUrl)
	if err != nil {
		return err
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)

	return nil
}

func createComment(ctx *gin.Context) {
	var cmt Comment
	if err := ctx.BindJSON(&cmt); err != nil {
		ctx.AbortWithStatusJSON(400, gin.H{
			"success": false,
			"message": err,
		})
	}

	cmtInBytes, err := json.Marshal(cmt)
	if err != nil {
		log.Fatal(err)
	}

	PushCommentToQueue("comments", cmtInBytes)

	ctx.JSON(200, gin.H{
		"success": true,
		"message": "Comment pushed successfully",
		"comment": cmt,
	})
}
