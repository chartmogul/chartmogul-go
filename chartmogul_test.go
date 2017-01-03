package chartmogul

import (
	"flag"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
)

var (
	cm    = flag.Bool("cm", false, "run integration library tests against ChartMogul")
	key   = flag.String("key", "", "access key for CM test")
	token = flag.String("token", "", "access token for CM test")
	api   = API{}
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Verbose() {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if *key == "" || *token == "" {
		if *cm {
			logrus.Info("Please supply testing account key and token on cmd line to run live tests.")
		}
		*cm = false
	}

	result := m.Run()

	os.Exit(result)
}

func TestPing(t *testing.T) {
	if !*cm {
		t.SkipNow()
		return
	}
	b, err := api.Ping()
	if err != nil {
		t.Error(err)
	} else if !b {
		t.Error("ping returned false")
	}
}

//TODO: unit tests against mocked HTTP server.
