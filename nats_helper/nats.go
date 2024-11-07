package nats_helper

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func NewNatsJetStream(nc *nats.EncodedConn, logger *zap.SugaredLogger) (jetstream.JetStream, error) {
	js, err := jetstream.New(nc.Conn)
	if err != nil {
		logger.Error("failed to create nats jetstream", zap.Error(err))
		return nil, err
	}

	return js, nil
}


func NewNatsConn(natsUrl string, logger *zap.SugaredLogger) (*nats.EncodedConn, error) {
    nc, err := nats.Connect(
		natsUrl,
		nats.RetryOnFailedConnect(true),
		//nats.MaxReconnects(100),
		nats.PingInterval(time.Second*30),
		nats.ReconnectWait(time.Second),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			logger.Infof("attempting to connect to nats server %s", natsUrl)
		}),
		nats.DisconnectErrHandler(func(c *nats.Conn, error error) {
			logger.Errorf("disconnected from nats %v", error)
			return
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			logger.Errorf("connection closed")
			return
		}))

    if err != nil {
		logger.Errorf("failed to connect to nats server %s: %v", natsUrl, err)
		return nil, err
	}

	enc, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		logger.Errorf("failed to create json encoder for nats client: %v", err)
		return nil, err
	}

	return enc, nil
}