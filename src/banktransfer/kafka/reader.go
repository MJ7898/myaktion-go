package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MJ7898/myaktion-go/src/banktransfer/grpc/banktransfer"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type TransactionReader interface {
	Close() error
	Read(ctx context.Context, handler func(transaction *banktransfer.Transaction) error) error
}

type ReadHandler func (transaction *banktransfer.Transaction) error

type reader struct {
	reader *kafka.Reader
}

func NewTransactionReader() *reader {
	return &reader{reader: kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{connect},
		GroupID: groupID,
		Topic: Topic,
		MinBytes: 10e2, // 1KB
		MaxBytes: 10e5, // 1MB
	})}
}

func (r *reader) Read(ctx context.Context, handler ReadHandler) error  {
	msg, err := r.reader.FetchMessage(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch message from Kafka server: %w", err)
	}
	entry := log.WithField("message", msg)
	entry.Trace("retrieve message from Kafka server")
	var transaction banktransfer.Transaction
	if err := json.Unmarshal(msg.Value, &transaction); err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}
	err = handler(&transaction)
	if err != nil {
		return err
	}
	if err := r.reader.CommitMessages(ctx, msg); err != nil {
		return fmt.Errorf("failed to commit message: %w", err)
	}
	return nil
}

func (r *reader) Close() error {
	if err := r.reader.Close(); err != nil {
		return fmt.Errorf("could not close Kafka subscription: %w", err)
	}
	return nil
}
