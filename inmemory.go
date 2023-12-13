package noder

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"sync"
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

	fm sync.Mutex
	hm sync.Mutex
}

func NewInmemoryNodesStorage() NodesStorage {
	return &nodesStorage{
		full:    []Node{},
		history: []Node{},
	}
}

func (s *nodesStorage) Add(node Node) error {
	if node.History {
		s.hm.Lock()
		s.history = append(s.history, node)
		s.hm.Unlock()
	} else {
		s.fm.Lock()
		s.full = append(s.full, node)
		s.fm.Unlock()
	}

	return nil
}

func (s *nodesStorage) Get() (*Node, error) {
	if len(s.full) == 0 {
		return nil, ErrNoHealthyNodesFound
	}

	if node := s.currentF(); node.CheckHealth() {
		return &node, nil
	}

	s.fm.Lock()

	// ensure currentF node was not changed by another goroutine
	if node := s.currentF(); node.CheckHealth() {
		s.fm.Unlock()
		return &node, nil
	}

	// circular slice iteration without checking currentF unhealthy node
	for i := 0; i < len(s.full)-1; i++ {
		idx := (s.fullCurrIdx + i + 1) % len(s.full)
		if node := s.full[idx]; node.CheckHealth() {
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

	if node := s.currentH(); node.CheckHealth() {
		return &node, nil
	}

	s.hm.Lock()

	// ensure currentH node was not changed by another goroutine
	if node := s.currentH(); node.CheckHealth() {
		s.hm.Unlock()
		return &node, nil
	}

	// circular slice iteration without checking currentH unhealthy node
	for i := 0; i < len(s.history)-1; i++ {
		idx := (s.historyCurrIdx + i + 1) % len(s.history)
		if node := s.history[idx]; node.CheckHealth() {
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
		if node := s.nextF(); node.CheckHealth() {
			return &node, nil
		}
	}

	return nil, ErrNoHealthyNodesFound
}

func (s *nodesStorage) GetNextHistory() (*Node, error) {
	s.fm.Lock()
	defer s.fm.Unlock()

	for i := 0; i < len(s.history); i++ {
		if node := s.nextH(); node.CheckHealth() {
			return &node, nil
		}
	}

	return nil, ErrNoHealthyNodesFound
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
