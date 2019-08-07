package p2p

import (
	libp2p_peer "github.com/libp2p/go-libp2p-peer"
	ma "github.com/multiformats/go-multiaddr"
)

type Peer struct {
	IP            string
	Port          int
	TargetAddress []ma.Multiaddr
	PeerID        libp2p_peer.ID
}
