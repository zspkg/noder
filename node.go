package noder

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Node struct {
	Rpc    *ethclient.Client `fig:"rpc,required"`
	RpcUrl string            `fig:"rpc,required"`

	Ws    *ethclient.Client `fig:"ws"`
	WsUrl string            `fig:"ws"`

	ChainTitle string `fig:"chain"`
	ChainId    int64  `fig:"chain_id,required"`
}

// CheckHealth tries to get current block number and
// if it fails, returns false (and true otherwise, respectively)
func (n Node) CheckHealth() bool {
	_, err := n.Rpc.BlockNumber(context.Background())
	return err == nil
}

// Name returns a string representation of the node
func (n Node) Name() string {
	return fmt.Sprintf("%s | %v", n.RpcUrl, n.ChainId)
}
