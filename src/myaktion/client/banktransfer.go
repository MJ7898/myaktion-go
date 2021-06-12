package client

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"os"
	"time"
)

var (
	TimeoutLimitation  = time.Second * 10
	bankTransferTarget = os.Getenv("BANKTRANSFER_CONNECT")
)

func GetBankTransferConnection() (*grpc.ClientConn, error) {
	var err error
	log.WithFields(log.Fields{
		"target": bankTransferTarget,
	}).Infoln("Connecting to banktransfer service")
	var conn *grpc.ClientConn
	ctx, _ := context.WithTimeout(context.Background(), TimeoutLimitation)
	conn, err = grpc.DialContext(ctx, bankTransferTarget, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
