package noder

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"sync"
)

var (
	ErrorNoSpecifiedNodesFound = errors.New("no specified nodes found")
)

type imcs struct {
	storages map[int64]NodesStorage
	m        sync.Mutex
}

func NewInmemoryMultiChainNodesStorage() MultiChainNodesStorage {
	return &imcs{storages: map[int64]NodesStorage{}}
}

func (s *imcs) Add(node Node) error {
	s.m.Lock()
	defer s.m.Unlock()

	if storage, ok := s.storages[node.ChainId]; !ok {
		s.storages[node.ChainId] = &nodesStorage{Nodes: []Node{node}}
	} else {
		_ = storage.Add(node)
	}

	return nil
}

func (s *imcs) GetByChainId(chainId int64) (*Node, error) {
	storage, ok := s.storages[chainId]
	if !ok {
		return nil, nil
	}

	return storage.Get()
}

func (s *imcs) ToSingleChain(chainId int64) (NodesStorage, error) {
	storage, ok := s.storages[chainId]
	if !ok {
		return nil, errors.From(ErrorNoSpecifiedNodesFound, logan.F{"chainId": chainId})
	}

	return storage, nil
}
