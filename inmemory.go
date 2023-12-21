package noder

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"sync"
	"time"
)

var (
	ErrNoHealthyNodesFound = errors.New("healthy node was not found")
)

type nodesStorage struct {
	full        []Node
	fullCurrIdx int
	fullNextIdx int

	history        []Node
	historyCurrIdx int
	historyNextIdx int

	healthyMap map[string]bool

	fm       sync.RWMutex
	hm       sync.RWMutex
	healthyM sync.RWMutex
}

func NewInmemoryNodesStorage(nodes []Node, healthCheckPeriod time.Duration) NodesStorage {
	storage := &nodesStorage{
		full:       []Node{},
		history:    []Node{},
		healthyMap: map[string]bool{},
	}

	for _, node := range nodes {
		if node.History {
			storage.history = append(storage.history, node)
		} else {
			storage.full = append(storage.full, node)
		}
	}

	go storage.watchHealth(healthCheckPeriod)

	return storage
}

func (s *nodesStorage) Get() (*Node, error) {
	if len(s.full) == 0 {
		return nil, ErrNoHealthyNodesFound
	}

	if node := s.currentF(); s.healthy(node) {
		return &node, nil
	}

	s.fm.Lock()

	// ensure currentF node was not changed by another goroutine
	if node := s.currentF(); s.healthy(node) {
		s.fm.Unlock()
		return &node, nil
	}

	// circular slice iteration without checking currentF unhealthy node
	for i := 0; i < len(s.full)-1; i++ {
		idx := (s.fullCurrIdx + i + 1) % len(s.full)
		if node := s.full[idx]; s.healthy(node) {
			s.fullCurrIdx = idx
			s.fm.Unlock()
			return &node, nil
		}
	}

	s.fm.Unlock()

	return s.GetHistory()
}

func (s *nodesStorage) GetHistory() (*Node, error) {
	if len(s.history) == 0 {
		return nil, ErrNoHealthyNodesFound
	}

	if node := s.currentH(); s.healthy(node) {
		return &node, nil
	}

	s.hm.Lock()

	// ensure currentH node was not changed by another goroutine
	if node := s.currentH(); s.healthy(node) {
		s.hm.Unlock()
		return &node, nil
	}

	// circular slice iteration without checking currentH unhealthy node
	for i := 0; i < len(s.history)-1; i++ {
		idx := (s.historyCurrIdx + i + 1) % len(s.history)
		if node := s.history[idx]; s.healthy(node) {
			s.historyCurrIdx = idx
			s.hm.Unlock()
			return &node, nil
		}
	}

	s.hm.Unlock()
	return nil, ErrNoHealthyNodesFound
}

func (s *nodesStorage) GetNext() (*Node, error) {
	s.hm.Lock()
	defer s.hm.Unlock()

	for i := 0; i < len(s.full); i++ {
		if node := s.nextF(); s.healthy(node) {
			return &node, nil
		}
	}

	return nil, ErrNoHealthyNodesFound
}

func (s *nodesStorage) GetNextHistory() (*Node, error) {
	s.fm.Lock()
	defer s.fm.Unlock()

	for i := 0; i < len(s.history); i++ {
		if node := s.nextH(); s.healthy(node) {
			return &node, nil
		}
	}

	return nil, ErrNoHealthyNodesFound
}

func (s *nodesStorage) watchHealth(period time.Duration) {
	ticker := time.NewTicker(period)

	for {
		for i := 0; i < len(s.full); i++ {
			s.healthyM.Lock()
			s.healthyMap[s.full[i].Name()] = s.full[i].CheckHealth()
			s.healthyM.Unlock()
		}
		for i := 0; i < len(s.history); i++ {
			s.healthyM.Lock()
			s.healthyMap[s.history[i].Name()] = s.history[i].CheckHealth()
			s.healthyM.Unlock()
		}
		<-ticker.C
	}
}

func (s *nodesStorage) healthy(node Node) bool {
	s.healthyM.RLock()
	defer s.healthyM.RUnlock()

	return s.healthyMap[node.Name()]
}

func (s *nodesStorage) currentF() Node {
	return s.full[s.fullCurrIdx]
}

func (s *nodesStorage) currentH() Node {
	return s.history[s.historyCurrIdx]
}

func (s *nodesStorage) nextF() Node {
	node := s.full[s.fullNextIdx]
	s.fullNextIdx = (s.fullNextIdx + 1) % len(s.full)
	return node
}

func (s *nodesStorage) nextH() Node {
	node := s.history[s.historyNextIdx]
	s.historyNextIdx = (s.historyNextIdx + 1) % len(s.history)
	return node
}
