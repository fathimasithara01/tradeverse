package tests

import (
	"testing"

	signal "github.com/fathimasithara01/tradeverse/internal/admin/service/signal"
)

func TestCreateSignal(t *testing.T) {
	svc := signal.NewSignalService()

	sig := signal.Signal{
		TraderID: 1,
		Entry:    100.5,
		Target:   110.0,
		StopLoss: 95.0,
	}

	err := svc.Publish(sig)
	if err != nil {
		t.Fatalf("failed to publish signal: %v", err)
	}
}
