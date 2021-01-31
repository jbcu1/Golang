package kafka_client

import (
	"PTwork/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"sync"
	"time"
)

const(
	topic = "pdns_test"
	broker1Address = "10.0.180.179:9092"
	broker2Address = "10.0.180.182:9092"
	broker3Address = "10.0.180.183:9092"
)



func Produce(hosts <-chan string, ctx context.Context, wg *sync.WaitGroup,kafkaMessage models.KafkaMessage){


	defer func() {

		if r:=recover(); r!=nil{

			fmt.Println("Panic happened",r)

		}
	}()

	defer wg.Done()

	//l := log.New(os.Stdout, "kafka writer: ", 0)

	w:=kafka.NewWriter(kafka.WriterConfig{

		Brokers: []string{broker1Address,broker2Address,broker3Address},
		Topic: topic,
		//Logger: l,
		BatchSize: 10,
		BatchTimeout: 2 * time.Second,
		RequiredAcks: 1,

	})

	for host:=range hosts {

		reqBodyBytes:=new(bytes.Buffer)
		_ = json.NewEncoder(reqBodyBytes).Encode(kafkaMessage)
		val:=reqBodyBytes.Bytes()
		val=val[:len(val)-1]

		err:=w.WriteMessages(ctx,kafka.Message{

			Key:   []byte(host),
			Value: val,

		})

		if err!=nil{

			panic("could not write message "+ err.Error())

		}
	}
}

func Consumer(ctx context.Context,hosts chan<- string,tags chan<- string){

	defer func() {

		if r:=recover(); r!=nil{

			fmt.Println("Panic happened in consumer",r)

		}

	}()

	//l := log.New(os.Stdout, "kafka writer: ", 0)

	r:=kafka.NewReader(kafka.ReaderConfig{

		Brokers: []string{broker1Address,broker2Address,broker3Address},
		Topic: topic,
		GroupID: "pdns_cons",
		//Logger: l,
		StartOffset: kafka.FirstOffset,
		MaxWait: 1*time.Second,
		CommitInterval: time.Second,


	})





	for {

		m, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}

		hosts<-string(m.Key)
		tags<-string(m.Value)

		if err := r.CommitMessages(ctx, m); err != nil {

			log.Fatal("failed to commit messages:", err)

		}


	}

	r.Close()
	close(hosts)
	close(tags)
}