package noder

type Noder interface {
	MultiChainNodesStorage() MultiChainNodesStorage
	NodesStorage() NodesStorage
}

type MultiChainNodesStorage interface {
	Add(node Node) error
	GetByChainId(chainId int64) (*Node, error)
	ToSingleChain(chainId int64) (NodesStorage, error)
}

type NodesStorage interface {
	Add(node Node) error
	Get() (*Node, error)
}
