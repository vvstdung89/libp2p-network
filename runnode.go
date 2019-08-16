package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/multiformats/go-multiaddr"
	"log"
	"node/p2p"
	"time"
)

func main() {
	var port = flag.Int("port", 10001, "Use port")
	flag.Parse()
	node := p2p.NewNode(p2p.NodeConfig{Port: *port, PublicIP: "127.0.0.1", Version: "1.1", MaxConnection: 100})
	fmt.Println(node.Host.ID(), node.Host.Addrs())

	node.Subscribe("thethao", func(st *pubsub.Subscription) {
		for {
			m, e := st.Next(context.Background())
			if e != nil {
				log.Println(e)
			}
			log.Println("st1", string(m.Data))
		}
	})

	go func() {
		if *port == 10001 {

			addr1, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10000")
			addr2, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10002")
			addr3, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10003")
			addr4, _ := multiaddr.NewMultiaddr("/ip4/127.0.0.1/tcp/10004")
			id1, _ := peer.IDB58Decode("QmaPNiPsv61e6TSDrXM6GX1n9iLsBXLbkHe3xz3D4mcpKU")
			id2, _ := peer.IDB58Decode("QmaB3zjwcFSwcAsf2uAfX6rkwFQtzmA8Y7Gm129gXbj6QQ")
			id3, _ := peer.IDB58Decode("QmNton3rrwD9ifvHuMFySBigKWMwjhxoWRBMRUfNm4rikU")
			id4, _ := peer.IDB58Decode("QmSv8EABiz7LeDTEuiQ25KXzxMq7eSP51M5KeeQLumQeww")
			//
			if err := node.Host.Connect(context.Background(), peer.AddrInfo{id1, []multiaddr.Multiaddr{addr1}}); err != nil {
				fmt.Println("Cannot connect to 10000", err)
			}

			if err := node.Host.Connect(context.Background(), peer.AddrInfo{id2, []multiaddr.Multiaddr{addr2}}); err != nil {
				fmt.Println("Cannot connect to 10002", err)
			}
			if err := node.Host.Connect(context.Background(), peer.AddrInfo{id3, []multiaddr.Multiaddr{addr3}}); err != nil {
				fmt.Println("Cannot connect to 10003", err)
			}
			if err := node.Host.Connect(context.Background(), peer.AddrInfo{id4, []multiaddr.Multiaddr{addr4}}); err != nil {
				fmt.Println("Cannot connect to 10004", err)
			}

			a := time.Tick(1 * time.Second)
			for _ = range a {
				node.Host.ConnManager().TagPeer(id1, "sd", 4)
				fmt.Println(node.Host.ConnManager().GetTagInfo(id1))
			}

		}

	}()
	select {}
}
