package tests

import (
	"testing"
	"time"

	subscription "github.com/fathimasithara01/tradeverse/internal/admin/service/subscription"
)

func TestSubscriptionExpiry(t *testing.T) {
	s := subscription.NewSubscription(time.Now().Add(-1 * time.Hour))

	if !s.IsExpired() {
		t.Errorf("expected subscription to be expired")
	}
}
