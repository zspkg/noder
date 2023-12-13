package noder

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"time"
)

const healthyNodeTimeout = 2 * time.Second

type Node struct {
	Rpc    *ethclient.Client `fig:"rpc,required"`
	RpcUrl string            `fig:"rpc,required"`

	Ws    *ethclient.Client `fig:"ws"`
	WsUrl string            `fig:"ws"`

	ChainId int64 `fig:"chain_id,required"`

	History bool `fig:"history"`
}

// CheckHealth tries to get currentF block number and
// if it fails, returns false (and true otherwise, respectively)
func (n Node) CheckHealth() bool {
	ctx, _ := context.WithTimeout(context.Background(), healthyNodeTimeout)
	_, err := n.Rpc.BlockNumber(ctx)
	return err == nil
}

// Name returns a string representation of the node
func (n Node) Name() string {
	return fmt.Sprintf("%s | %v", n.RpcUrl, n.ChainId)
}
