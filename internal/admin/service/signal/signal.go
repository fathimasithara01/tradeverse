package signal

import "fmt"

type Signal struct {
	TraderID int
	Entry    float64
	Target   float64
	StopLoss float64
}

type Service struct{}

func NewSignalService() *Service {
	return &Service{}
}

func (s *Service) Publish(sig Signal) error {
	if sig.TraderID == 0 {
		return fmt.Errorf("trader id required")
	}
	return nil
}
