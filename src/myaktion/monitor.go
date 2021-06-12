package main

import (

	"github.com/MJ7898/myaktion-go/src/myaktion/client"
	"github.com/MJ7898/myaktion-go/src/myaktion/client/banktransfer"
	"github.com/MJ7898/myaktion-go/src/myaktion/service"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

const reconnectTimeout = 120 * time.Second

func monitortransactions() {
	for{
		connectandmonitor()
		time.Sleep(time.Second)
	}
}
func connectandmonitor() {
	conn, err := client.GetBankTransferConnection()
	if err != nil {
		log.WithError(err).Fatal("error connecting to the banktransfer service")
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), reconnectTimeout)
	defer cancel()
	banktransferClient := banktransfer.NewBankTransferClient(conn)
	watcher, err := banktransferClient.ProcessTransactions(ctx)
	if err != nil {
		log.WithError(err).Fatal("error watching transactions")
	}
	log.Info("Successfully connected to banktransfer service for processing transactions")

	for {
		transaction, err := watcher.Recv()
		if err != nil {
			if _, deadline := ctx.Deadline(); deadline {
				log.Info("deadline reached. reconnect client")
				break
			}
			log.WithError(err).Error("error receiving transaction")
			continue

		}

		entry := log.WithField("transaction", transaction)
		entry.Info("Received transaction")
		err = service.MarkDonation(uint(transaction.DonationId))
		if err != nil {
			entry.WithError(err).Error("error changing donation status")
			continue
		}
		entry.Info("Sending processing response")
		err = watcher.Send(&banktransfer.ProcessingResponse{Id: transaction.Id})
		if err != nil {
			entry.WithError(err).Error("error sending processing response")
			continue
		}
		entry.Info("Processing response sent")
	}
}