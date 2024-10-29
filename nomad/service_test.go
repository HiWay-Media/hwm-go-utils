package nomad_test

import (
	"os"
	"testing"

	"github.com/HiWay-Media/hwm-go-utils/log"
	"github.com/HiWay-Media/hwm-go-utils/nomad"
)

func TestMain(m *testing.M) {
	if os.Getenv("APP_ENV") == "" {
		err := os.Setenv("APP_ENV", "test")
		if err != nil {
			panic("could not set test env")
		}
	}
	//env.Load()
	m.Run()
}

func getNomad() nomad.IService {
	options := nomad.Options{
		BaseUrl:  "url",
		LogLevel: "debug",
		Logger:   log.GetLogger("debug"),
	}
	return nomad.NewService(options)
}
