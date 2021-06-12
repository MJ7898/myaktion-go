package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/MJ7898/myaktion-go/src/banktransfer/grpc/banktransfer"
	"github.com/MJ7898/myaktion-go/src/banktransfer/kafka"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"math/rand"
	"time"
)

const retryTime = 5 * time.Second

type BankTransferService struct {
	banktransfer.BankTransferServer
	keyGenerator *KeyGenerator
	transactionWriter kafka.TransactionWriter
	// NOTE: should be atomic if we have multiple connections
}

func NewBankTransferService() *BankTransferService {
	rand.Seed(time.Now().UnixNano())
	return &BankTransferService{keyGenerator: NewKeyGenerator(),
	}
}

func (s *BankTransferService) TransferMoney(_ context.Context, transaction *banktransfer.Transaction) (*emptypb.Empty, error) {
	// log.Infof("Received transaction: %v", transaction)
	entry := log.WithField("transaction", transaction)
	entry.Info("Received transaction")
	s.processTransaction(transaction)
	return &emptypb.Empty{}, nil
}

func (s *BankTransferService) ProcessTransactions(stream banktransfer.BankTransfer_ProcessTransactionsServer) error {
	return func() error {
		r := kafka.NewTransactionReader()
		for {
			err := r.Read(stream.Context(), func(transaction *banktransfer.Transaction) error {
				id := transaction.Id
				entry := log.WithField("transaction", transaction)
				entry.Info("Sending transaction")
				if err := stream.Send(transaction); err != nil {
					return fmt.Errorf("error sending transaction: %w", err)
				}
				entry.Info("Transaction sent. Waiting for processing response")
				response, err := stream.Recv()
				if err != nil {
					return fmt.Errorf("error receiving processing response: %w", err)
				}
				if response.Id != id {
					// NOTE: this is just a guard and not happening as transaction is local per connection
					return errors.New("received processing response of a different transaction")
				}
				entry.Info("Processing response received")
				return nil
			})
			if err != nil {
				log.WithError(err).Error("error while reading transaction")
				break
			}
		}
		r.Close()
		return nil
	}()
}

func (s *BankTransferService) processTransaction(transaction *banktransfer.Transaction) {
	entry := log.WithField("transaction", transaction)
	go func(transaction banktransfer.Transaction) {
		entry.Info("Start processing transaction")
		transaction.Id = s.keyGenerator.getUniqueId() // atomic.AddInt32(&s.counter, 1)
		if err := s.transactionWriter.Write(&transaction); err != nil {
			entry.WithError(err).Error("Can't write transaction to transaction writer")
			return
		}
		entry.Info("Transaction forwarded to transaction writer. Processing transaction finished")
	}(*transaction)
}

func (s *BankTransferService) Start() {
	log.Info("Starting banktransfer service")
	s.transactionWriter = kafka.NewTransactionWriter()
	log.Info("Successfully created transaction writer")
}

func (s *BankTransferService) Stop() {
	log.Info("Stopping banktransfer service")
	s.transactionWriter.Close()
	log.Info("Successfully closed connection to transaction writer")
}
