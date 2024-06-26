package piko

import (
	"context"

	"github.com/andydunstall/piko/pkg/log"
)

const (
	// defaultURL is the URL of the Piko upstream port when running locally.
	defaultURL = "ws://localhost:8001"
)

// Piko manages registering and listening on endpoints.
//
// The client establishes an outbound-only connection to the server for each
// listener. Proxied connections for the listener are then multiplexed over
// that outbound connection. Therefore the client never exposes a port.
type Piko struct {
	options options
	logger  log.Logger
}

func New(opts ...Option) *Piko {
	options := options{
		token:  "",
		url:    defaultURL,
		logger: log.NewNopLogger(),
	}
	for _, o := range opts {
		o.apply(&options)
	}

	return &Piko{
		options: options,
		logger:  options.logger,
	}
}

// Listen listens for connections for the given endpoint ID.
//
// Listen will block until the listener has been registered.
//
// The returned [Listener] is a [net.Listener].
func (p *Piko) Listen(ctx context.Context, endpointID string) (Listener, error) {
	return listen(ctx, endpointID, p.options, p.logger)
}
