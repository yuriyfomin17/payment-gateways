package services

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

type FaultToleranceService struct {
	circuitBreaker *gobreaker.CircuitBreaker
}

func NewFaultToleranceService() FaultToleranceService {
	return FaultToleranceService{
		circuitBreaker: gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        "RabbitMqPublisher",
			MaxRequests: 1,
			Interval:    5 * time.Second,
			Timeout:     3 * time.Second,
		}),
	}
}

// PublishWithCircuitBreaker publishes a transaction to RabbitMQ using a circuit breaker
func (ft FaultToleranceService) PublishWithCircuitBreaker(operation func() error) error {
	_, err := ft.circuitBreaker.Execute(func() (any, error) {
		return nil, operation()
	})
	return err
}

// RetryOperation Retry operation with exponential backoff
func (ft FaultToleranceService) RetryOperation(operation func() error, maxRetries int) error {
	for i := 0; i < maxRetries; i++ {
		if err := operation(); err == nil {
			return nil
		}
		time.Sleep(time.Duration(2^i) * time.Second)
	}
	return fmt.Errorf("operation failed after %d attempts", maxRetries)
}
