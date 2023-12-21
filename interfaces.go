package noder

type Noder interface {
	NodesConfig() NodesConfig
	MultiChainNodesStorage() MultiChainNodesStorage
	NodesStorage() NodesStorage
}

type MultiChainNodesStorage interface {
	GetByChainId(chainId int64) (*Node, error)
	GetNextByChainId(chainId int64) (*Node, error)
	GetHistoryByChainId(chainId int64) (*Node, error)
	GetNextHistoryByChainId(chainId int64) (*Node, error)
	ToSingleChain(chainId int64) (NodesStorage, error)
}

type NodesStorage interface {
	Get() (*Node, error)
	GetNext() (*Node, error)
	GetHistory() (*Node, error)
	GetNextHistory() (*Node, error)
}
