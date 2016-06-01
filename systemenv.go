package septum

import (
	"sync/atomic"
	"time"
)

type (
	systemEnvironment struct {
		nextEventId uint64
	}
)

// ==== Environment ====

var (
	systemEnv Environment
)

func init() {
	systemEnv = NewSystemEnvironment(0)
}

func SystemEnvironment() Environment {
	return systemEnv
}

func NewSystemEnvironment(eventIdSeed uint64) Environment {
	return &systemEnvironment{
		nextEventId: eventIdSeed,
	}
}

func (e *systemEnvironment) NextEventId() uint64 {
	return atomic.AddUint64(&e.nextEventId, 1)
}

func (e *systemEnvironment) Now() time.Time {
	return time.Now()
}
