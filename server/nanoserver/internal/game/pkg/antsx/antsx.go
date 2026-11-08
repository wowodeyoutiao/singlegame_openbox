package antsx

import (
	"log"
	"time"

	"github.com/lonng/nano/scheduler"
	"github.com/panjf2000/ants/v2"
)

// LocalScheduler is the concrete implementation of the LocalScheduler interface.
// It uses an ants.Pool to manage a pool of goroutines.
type LocalScheduler struct {
	pool *ants.Pool
}

// options holds the configurable parameters for the scheduler's ants pool.
type options struct {
	// Size is the maximum number of goroutines in the pool.
	Size int
	// ExpiryDuration is the period for which an idle goroutine is kept alive.
	ExpiryDuration time.Duration
	// PreAlloc indicates if memory for goroutines should be pre-allocated.
	PreAlloc bool
	// Logger is the custom logger for the ants pool.
	Logger ants.Logger
	// PanicHandler is the function to handle panics from tasks.
	PanicHandler func(interface{})
}

// Option is a function that configures the scheduler's options.
type Option func(*options)

// WithSize sets the maximum number of goroutines in the pool.
// The default is ants.DefaultAntsPoolSize.
func WithSize(size int) Option {
	return func(o *options) {
		o.Size = size
	}
}

// WithExpiryDuration sets the duration for which an idle goroutine is kept alive.
// The default is 10 seconds.
func WithExpiryDuration(duration time.Duration) Option {
	return func(o *options) {
		o.ExpiryDuration = duration
	}
}

// WithPreAlloc enables or disables pre-allocation of memory for goroutines.
func WithPreAlloc(preAlloc bool) Option {
	return func(o *options) {
		o.PreAlloc = preAlloc
	}
}

// WithLogger sets a custom logger for the pool.
// The default is a standard log.Logger.
func WithLogger(logger ants.Logger) Option {
	return func(o *options) {
		o.Logger = logger
	}
}

// WithPanicHandler sets a function to handle panics from tasks.
// If not set, panics will be propagated.
func WithPanicHandler(handler func(interface{})) Option {
	return func(o *options) {
		o.PanicHandler = handler
	}
}

// NewLocalScheduler creates and returns a new LocalScheduler instance.
// It initializes an ants pool based on the provided configuration options.
func NewLocalScheduler(opts ...Option) (scheduler.LocalScheduler, error) {
	// 1. Define default options.
	defaultOpts := &options{
		Size:           ants.DefaultAntsPoolSize,
		ExpiryDuration: 10 * time.Second,
		PreAlloc:       false,
		Logger:         log.Default(),
	}

	// 2. Apply user-provided options to override defaults.
	for _, opt := range opts {
		opt(defaultOpts)
	}

	// 3. Create the ants pool with the configured options.
	pool, err := ants.NewPool(
		defaultOpts.Size,
		ants.WithExpiryDuration(defaultOpts.ExpiryDuration),
		ants.WithPreAlloc(defaultOpts.PreAlloc),
		ants.WithLogger(defaultOpts.Logger),
		ants.WithPanicHandler(defaultOpts.PanicHandler),
	)
	if err != nil {
		return nil, err
	}

	return &LocalScheduler{pool: pool}, nil
}

// Schedule submits a task to the ants goroutine pool for execution.
// It adheres to the LocalScheduler interface.
func (s *LocalScheduler) Schedule(task scheduler.Task) {
	err := s.pool.Submit(task)
	if err != nil {
		log.Printf("Failed to submit task to scheduler: %v", err)
	}
}

// Release gracefully stops the goroutine pool.
func (s *LocalScheduler) Release() {
	s.pool.Release()
}
