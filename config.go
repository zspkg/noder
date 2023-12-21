package noder

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
	"reflect"
	"time"
)

const defaultConfigKey = "nodes"

type in struct {
	nodes    comfig.Once
	storage  comfig.Once
	mstorage comfig.Once
	getter   kv.Getter
	kvKey    string
}

type NodesConfig struct {
	Nodes             []Node        `fig:"nodes,required"`
	HealthCheckPeriod time.Duration `fig:"health_check_period"`
}

func NewInmemoryHealthy(getter kv.Getter, kvKey *string) Noder {
	n := in{getter: getter, kvKey: defaultConfigKey}
	if kvKey != nil {
		n.kvKey = *kvKey
	}

	return &n
}

func (n *in) MultiChainNodesStorage() MultiChainNodesStorage {
	return n.mstorage.Do(func() interface{} {
		return NewInmemoryMultiChainNodesStorage(n.NodesConfig())
	}).(MultiChainNodesStorage)
}

func (n *in) NodesStorage() NodesStorage {
	return n.storage.Do(func() interface{} {
		cfg := n.NodesConfig()
		return NewInmemoryNodesStorage(cfg.Nodes, cfg.HealthCheckPeriod)
	}).(NodesStorage)
}

func (n *in) NodesConfig() NodesConfig {
	return n.nodes.Do(func() interface{} {
		nodesCfg := NodesConfig{}

		if err := figure.
			Out(&nodesCfg).
			With(figure.BaseHooks, nodesHooks, evmHooks).
			From(kv.MustGetStringMap(n.getter, n.kvKey)).
			Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out nodes"))
		}

		if len(nodesCfg.Nodes) == 0 {
			panic(errors.New("no nodes were provided"))
		}

		return nodesCfg
	}).(NodesConfig)
}

var nodesHooks = figure.Hooks{
	"[]noder.Node": func(value any) (reflect.Value, error) {
		switch s := value.(type) {
		case []any:
			nodes := make([]Node, 0, len(s))
			for _, raw := range s {
				var node Node
				stringMap := make(map[string]any)

				switch v := raw.(type) {
				case map[string]any:
					stringMap = v
				case map[any]any:
					for k, v := range v {
						if str, ok := k.(string); !ok {
							return reflect.Value{}, errors.New("failed to cast to map[string]any")
						} else {
							stringMap[str] = v
						}
					}
				default:
					return reflect.Value{}, errors.New("got wrong type while figuring out []noder.Node")
				}

				if err := figure.
					Out(&node).
					With(figure.BaseHooks, evmHooks).
					From(stringMap).
					Please(); err != nil {
					return reflect.Value{}, errors.Wrap(err, "failed to figure out")
				}

				nodes = append(nodes, node)
			}

			return reflect.ValueOf(nodes), nil
		default:
			return reflect.Value{}, errors.New("unexpected type while figuring out []noder.Node")
		}
	},
}

var evmHooks = figure.Hooks{
	"*ethclient.Client": func(value any) (reflect.Value, error) {
		switch v := value.(type) {
		case string:
			client, err := ethclient.Dial(v)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "failed to convert value into ethclient")
			}
			return reflect.ValueOf(client), nil
		default:
			return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
		}
	},
}
