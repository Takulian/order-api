package breaker

import (
	"errors"
	"sync"
	"time"
)

type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

type Settings struct {
	Name                string
	FailureThreshold    int
	Openduration        time.Duration
	HalfOpenMaxRequests int
	OnStateChange       func(name string, from, to State)
}

type CircuitBreaker struct {
	mu sync.Mutex

	name                string
	failureThreshold    int
	openDuration        time.Duration
	halfOpenMaxRequests int
	onStateCHange       func(name string, from, to State)

	state              State
	consecutiveFailure int
	openUntil          time.Time
	halfOpenFlight     int
}

func NewCircuitBreaker(settings Settings) *CircuitBreaker {
	if settings.FailureThreshold <= 0 {
		settings.FailureThreshold = 5
	}
	if settings.HalfOpenMaxRequests <= 0 {
		settings.HalfOpenMaxRequests = 2
	}
	if settings.Openduration <= 0 {
		settings.Openduration = 10 * time.Second
	}

	return &CircuitBreaker{
		name:                settings.Name,
		failureThreshold:    settings.FailureThreshold,
		openDuration:        settings.Openduration,
		halfOpenMaxRequests: settings.HalfOpenMaxRequests,
		onStateCHange:       settings.OnStateChange,
		state:               StateClosed,
	}
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	err := fn()
	cb.afterRequest(err)

	return err
}

func (cb *CircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateOpen:
		if time.Now().Before(cb.openUntil) {
			return ErrCircuitOpen
		}
		cb.transitionTo(StateHalfOpen)
		cb.halfOpenFlight = 0
	case StateHalfOpen:
		if cb.halfOpenFlight >= cb.halfOpenMaxRequests {
			return ErrCircuitOpen
		}
	}

	if cb.state == StateHalfOpen {
		cb.halfOpenFlight++
	}

	return nil
}

func (cb *CircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if err != nil {
		cb.onFailure()
		return
	}

	cb.onSuccess()
}

func (cb *CircuitBreaker) onFailure() {
	switch cb.state {
	case StateHalfOpen:
		cb.transitionTo(StateOpen)
		cb.openUntil = time.Now().Add(cb.openDuration)
	case StateClosed:
		cb.consecutiveFailure++
		if cb.consecutiveFailure >= cb.failureThreshold {
			cb.transitionTo(StateOpen)
			cb.openUntil = time.Now().Add(cb.openDuration)
		}
	}

}
func (cb *CircuitBreaker) onSuccess() {
	switch cb.state {
	case StateHalfOpen:
		cb.transitionTo(StateClosed)
		cb.consecutiveFailure = 0
	case StateClosed:
		cb.consecutiveFailure = 0
	}

}

func (cb *CircuitBreaker) transitionTo(newState State) {
	if cb.state == newState {
		return
	}
	oldState := cb.state
	cb.state = newState
	if cb.onStateCHange != nil {
		cb.onStateCHange(cb.name, oldState, newState)
	}
}
