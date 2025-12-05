package subscription

import "time"

type Subscription struct {
	EndDate time.Time
}

func NewSubscription(end time.Time) *Subscription {
	return &Subscription{EndDate: end}
}

func (s *Subscription) IsExpired() bool {
	return time.Now().After(s.EndDate)
}
