package corehttp

import (
	"fmt"
	"net/http"

	key "github.com/ipfs/go-ipfs/blocks/key"
	core "github.com/ipfs/go-ipfs/core"
	id "github.com/ipfs/go-ipfs/p2p/protocol/identify"
)

// Gateway should be instantiated using NewGateway
type Gateway struct {
	Config GatewayConfig
}

type GatewayConfig struct {
	BlackList key.KeySet
	WhiteList key.KeySet

	Writable bool
}

func NewGateway(conf GatewayConfig) *Gateway {
	return &Gateway{
		Config: conf,
	}
}

func (g *Gateway) ServeOption() ServeOption {
	return func(n *core.IpfsNode, mux *http.ServeMux) (*http.ServeMux, error) {
		gateway, err := newGatewayHandler(n, g.Config)
		if err != nil {
			return nil, err
		}
		mux.Handle("/ipfs/", gateway)
		mux.Handle("/ipns/", gateway)
		return mux, nil
	}
}

func GatewayOption(writable bool) ServeOption {
	g := NewGateway(GatewayConfig{
		Writable: writable,
	})
	return g.ServeOption()
}

func VersionOption() ServeOption {
	return func(n *core.IpfsNode, mux *http.ServeMux) (*http.ServeMux, error) {
		mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Client Version:   %s\n", id.ClientVersion)
			fmt.Fprintf(w, "Protocol Version: %s\n", id.IpfsVersion)
		})
		return mux, nil
	}
}
