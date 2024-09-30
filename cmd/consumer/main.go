package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/trace"
	"sync/atomic"
	"time"

	"github.com/bretanac93/ckg/internal/db"
	"github.com/bretanac93/ckg/internal/order"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const appID = "consumer-kafka-go-test"

var messageCount uint64

func main() {
	ctx := context.Background()

	// pprof server
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if err := trace.Start(f); err != nil {
		panic(err)
	}
	defer trace.Stop()

	// define topics to listen from. Also wildcards are supported, something like "local.*"
	// that's unsafe tho. I prefer to declare all topics explicitly.
	topics := []string{"local.payments", "local.orders"}

	conn, closeConn, err := db.Open()
	if err != nil {
		panic(err)
	}
	defer closeConn()

	ordersProcessor := order.NewProcessor(conn)

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     "127.0.0.1:29092",
		"group.id":              appID,
		"broker.address.family": "v4",
		"session.timeout.ms":    10000,
		"heartbeat.interval.ms": 3000,
		"max.poll.interval.ms":  300000,
		"fetch.min.bytes":       1,
		"fetch.wait.max.ms":     500,
	})
	if err != nil {
		panic(err)
	}

	if err := c.SubscribeTopics(topics, nil); err != nil {
		panic(err)
	}

	slog.Info("listening for messages...")

	// goroutine to log the number of messages processed per second
	go func() {
		for {
			time.Sleep(1 * time.Second)
			count := atomic.SwapUint64(&messageCount, 0)
			slog.Info("Messages processed per second", "count", count)
		}
	}()

	for {
		start := time.Now()

		msg, err := c.ReadMessage(-1)
		if err == nil {
			start = time.Now()
			o, err := order.FromJSON(msg.Value)

			if err != nil {
				slog.Error("Error while unmarshalling order", "error", err.Error())
				continue
			}
			unmarshalDuration := time.Since(start)
			slog.Info("time taken to unmarshal message", "duration", unmarshalDuration)

			start = time.Now()

			if err := ordersProcessor.Process(ctx, o); err != nil {
				slog.Error("error while processing order", "error", err.Error())
				continue
			}

			processDuration := time.Since(start)
			slog.Info("time taken to process message", "duration", processDuration)

			atomic.AddUint64(&messageCount, 1)
		} else if !err.(kafka.Error).IsTimeout() {
			slog.Warn("the consumer has been stale for a while")
		}
	}
}
