package noder

import (
	"gitlab.com/distributed_lab/logan/v3/errors"
	"sync"
)

var (
	ErrNoHealthyNodesFound = errors.New("healthy node was not found")
)

type nodesStorage struct {
	Nodes            []Node
	CurrentNodeIndex int
	m                sync.Mutex
}

func NewInmemoryNodesStorage() NodesStorage {
	return &nodesStorage{
		Nodes: []Node{},
	}
}

func (s *nodesStorage) Add(node Node) error {
	s.m.Lock()
	defer s.m.Unlock()

	s.Nodes = append(s.Nodes, node)

	return nil
}

func (s *nodesStorage) current() Node {
	return s.Nodes[s.CurrentNodeIndex]
}

func (s *nodesStorage) Get() (*Node, error) {
	if len(s.Nodes) == 0 {
		return nil, ErrNoHealthyNodesFound
	}

	if node := s.current(); node.CheckHealth() {
		return &node, nil
	}

	s.m.Lock()
	defer s.m.Unlock()

	// ensure current node was not changed by another goroutine
	if node := s.current(); node.CheckHealth() {
		return &node, nil
	}

	// circular slice iteration without checking current unhealthy node
	for i := 0; i < len(s.Nodes)-1; i++ {
		idx := (s.CurrentNodeIndex + i + 1) % len(s.Nodes)
		if node := s.Nodes[idx]; node.CheckHealth() {
			s.CurrentNodeIndex = idx
			return &node, nil
		}
	}

	return nil, ErrNoHealthyNodesFound
}
