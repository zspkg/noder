package noder

import (
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

var (
	ErrorNoSpecifiedNodesFound = errors.New("no specified nodes found")
)

type imcs struct {
	storages map[int64]NodesStorage
}

func NewInmemoryMultiChainNodesStorage(cfg NodesConfig) MultiChainNodesStorage {
	nodesMap := map[int64][]Node{}
	for _, node := range cfg.Nodes {
		nodesMap[node.ChainId] = append(nodesMap[node.ChainId], node)
	}

	storage := &imcs{storages: map[int64]NodesStorage{}}
	for chainId, nodes := range nodesMap {
		storage.storages[chainId] = NewInmemoryNodesStorage(nodes, cfg.HealthCheckPeriod)
	}

	return storage
}

func (s *imcs) GetByChainId(chainId int64) (*Node, error) {
	storage, ok := s.storages[chainId]
	if !ok {
		return nil, ErrorNoSpecifiedNodesFound
	}

	return storage.Get()
}

func (s *imcs) GetNextByChainId(chainId int64) (*Node, error) {
	storage, ok := s.storages[chainId]
	if !ok {
		return nil, ErrorNoSpecifiedNodesFound
	}

	return storage.GetNext()
}

func (s *imcs) GetHistoryByChainId(chainId int64) (*Node, error) {
	storage, ok := s.storages[chainId]
	if !ok {
		return nil, ErrorNoSpecifiedNodesFound
	}

	return storage.GetHistory()
}

func (s *imcs) GetNextHistoryByChainId(chainId int64) (*Node, error) {
	storage, ok := s.storages[chainId]
	if !ok {
		return nil, ErrorNoSpecifiedNodesFound
	}

	return storage.GetNextHistory()
}

func (s *imcs) ToSingleChain(chainId int64) (NodesStorage, error) {
	storage, ok := s.storages[chainId]
	if !ok {
		return nil, errors.From(ErrorNoSpecifiedNodesFound, logan.F{"chainId": chainId})
	}

	return storage, nil
}
