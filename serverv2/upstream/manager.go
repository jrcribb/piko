package upstream

import (
	"fmt"
	"net"
	"sync"
)

// Manager manages the set of local upsteam services.
type Manager struct {
	upstreams map[string][]*Upstream

	mu sync.Mutex
}

func NewManager() *Manager {
	return &Manager{
		upstreams: make(map[string][]*Upstream),
	}
}

func (m *Manager) Add(upstream *Upstream) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.upstreams[upstream.EndpointID] = append(
		m.upstreams[upstream.EndpointID], upstream,
	)
}

func (m *Manager) Remove(upstream *Upstream) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var upstreams []*Upstream
	for _, u := range m.upstreams[upstream.EndpointID] {
		if u != upstream {
			upstreams = append(upstreams, u)
		}
	}

	m.upstreams[upstream.EndpointID] = upstreams
}

func (m *Manager) Dial(endpointID string) (net.Conn, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	upstreams, ok := m.upstreams[endpointID]
	if !ok || len(upstreams) == 0 {
		return nil, fmt.Errorf("not found")
	}

	stream, err := upstreams[0].Sess.OpenTypedStream(0)
	if err != nil {
		return nil, err
	}
	return stream, nil
}
