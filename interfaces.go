package noder

type Noder interface {
	Nodes() []Node
	MultiChainNodesStorage() MultiChainNodesStorage
	NodesStorage() NodesStorage
}

type MultiChainNodesStorage interface {
	Add(node Node) error
	GetByChainId(chainId int64) (*Node, error)
	GetNextByChainId(chainId int64) (*Node, error)
	GetHistoryByChainId(chainId int64) (*Node, error)
	GetNextHistoryByChainId(chainId int64) (*Node, error)
	ToSingleChain(chainId int64) (NodesStorage, error)
}

type NodesStorage interface {
	Add(node Node) error
	Get() (*Node, error)
	GetNext() (*Node, error)
	GetHistory() (*Node, error)
	GetNextHistory() (*Node, error)
}
