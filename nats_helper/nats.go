package nats_helper

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/zap"
)

func NewNatsJetStream(nc *nats.EncodedConn, logger *zap.SugaredLogger) (jetstream.JetStream, error) {
	js, err := jetstream.New(nc.Conn)
	if err != nil {
		logger.Errorf("failed to create nats jetstream: %v", err)
		return nil, err
	}

	return js, nil
}

func NewNatsConn(natsUrl string, logger *zap.SugaredLogger) (*nats.EncodedConn, error) {
	nc, err := nats.Connect(
		natsUrl,
		nats.RetryOnFailedConnect(true),
		// keep retrying forever: the default (60 attempts) permanently closes the
		// connection after about a minute of server downtime
		nats.MaxReconnects(-1),
		nats.PingInterval(time.Second*30),
		nats.ReconnectWait(time.Second),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			logger.Infof("reconnected to nats server %s", conn.ConnectedUrlRedacted())
		}),
		nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
			if err != nil {
				logger.Errorf("disconnected from nats: %v", err)
			}
		}),
		nats.ClosedHandler(func(c *nats.Conn) {
			logger.Errorf("nats connection closed")
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
