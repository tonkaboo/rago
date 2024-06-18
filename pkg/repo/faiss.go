package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"gems/pkg/chunk"
	"sync/atomic"
	"time"

	pb "gems/proto/faiss"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/streadway/amqp"
)

type faiss struct {
	mqConn         *amqp.Connection
	mqCh           *amqp.Channel
	sentChunkCount int32
	ackCount       int32
	indexBuilt     bool
	grpcClient     pb.FaissServiceClient
}

func NewFaiss() (*faiss, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		return nil, fmt.Errorf("failed to start a connection: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %s", err)
	}

	queues := [2]string{"chunks", "acks"}
	for _, q := range queues {
		if _, err := ch.QueueDeclare(
			q,     // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		); err != nil {
			return nil, fmt.Errorf("failed to declare a queue: %s", err)
		}
	}

	// Dial to grpc server
	grpcConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to grpc server: %s", err)
	}
	// Create the client using the connection
	client := pb.NewFaissServiceClient(grpcConn)

	return &faiss{mqConn: conn, mqCh: ch, grpcClient: client}, nil
}

func (f *faiss) InsertChunk(appChunk *chunk.Chunk) error {
	jsonChunk, err := json.Marshal(appChunk)
	if err != nil {
		return fmt.Errorf("failed to marshal the json chunk: %s", err)
	}

	err = f.mqCh.Publish(
		"",       // exchange
		"chunks", // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(jsonChunk),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish the message: %s", err)
	}
	fmt.Println(" [x] Sent a chunk")
	atomic.AddInt32(&f.sentChunkCount, 1)

	return nil
}

func (f *faiss) RetrieveRelatedChunks(questionEmbedding []float32) ([]string, error) {
	// ack计数和已发送的chunk数量一致，说明index已经构建完成
	// 否则继续等待，直到构建完成
	if !f.indexBuilt {
		if err := f.waitForIndex(); err != nil {
			return nil, fmt.Errorf("failed to check if the index is built: %s", err)
		}
	}

	// 调用rpc服务，发送questionEmbedding，接收返回的字符串切片
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &pb.SimilarityRequest{
		Embedding: questionEmbedding,
	}
	res, err := f.grpcClient.FindSimilarChunks(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar chunk: %s", err)
	}
	return res.ChunkContents, nil
}

func (f *faiss) CloseConn() {
	f.mqConn.Close()
	f.mqCh.Close()
}

func (f *faiss) CheckIndex() bool {
	return f.indexBuilt
}

func (f *faiss) waitForIndex() error {
	msgs, err := f.mqCh.Consume(
		"acks", // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to receive the ack messages: %s", err)
	}

	indexBuilt := make(chan bool, 1)
	go func() {
		defer close(indexBuilt)
		for range msgs {
			fmt.Println("[x] Received an ack")
			f.ackCount++
			if f.ackCount == f.sentChunkCount {
				f.indexBuilt = true
				indexBuilt <- true
				return
			}
		}
	}()
	<-indexBuilt
	// 应该加一个超时处理：万一永远收不到数量相当的ACK，那也得返回个错误啊
	return nil
}
