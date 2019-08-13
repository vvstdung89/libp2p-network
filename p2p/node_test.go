package p2p

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"testing"
	"time"
)

func Test_Connection(t *testing.T) {
	node1 := NewNode(NodeConfig{Port: 10000, PublicIP: "127.0.0.1", Version: "1.1", MaxConnection: 1})
	node2 := NewNode(NodeConfig{Port: 10001, PublicIP: "127.0.0.1", Version: "1.1", MaxConnection: 1})
	node3 := NewNode(NodeConfig{Port: 10001, PublicIP: "127.0.0.1", Version: "1.1", MaxConnection: 1})
	fmt.Println("node1", node1.Host.ID())
	fmt.Println("node2", node2.Host.ID())

	err := node2.Host.Connect(context.Background(), peer.AddrInfo{node1.Host.ID(), node1.Host.Addrs()})
	if err != nil {
		log.Println(err)
	}

	//time.Sleep(5 * time.Second)
	err = node3.Host.Connect(context.Background(), peer.AddrInfo{node1.Host.ID(), node1.Host.Addrs()})
	if err != nil {
		log.Println(err)
	}
	time.Sleep(1 * time.Second)
	fmt.Println(node1.Host.Network().Conns())
	fmt.Println(node2.Host.Network().Conns())
	fmt.Println(node3.Host.Network().Conns())
}

func Test_Pubsub(t *testing.T) {
	node1 := NewNode(NodeConfig{Port: 10000, PublicIP: "127.0.0.1", Version: "1.1", MaxConnection: 10})
	node2 := NewNode(NodeConfig{Port: 10001, PublicIP: "127.0.0.1", Version: "1.1", MaxConnection: 10})
	node3 := NewNode(NodeConfig{Port: 10002, PublicIP: "127.0.0.1", Version: "1.1", MaxConnection: 10})

	err := node2.Host.Connect(context.Background(), peer.AddrInfo{node1.Host.ID(), node1.Host.Addrs()})
	if err != nil {
		log.Println(err)
	}

	err = node3.Host.Connect(context.Background(), peer.AddrInfo{node1.Host.ID(), node1.Host.Addrs()})
	if err != nil {
		log.Println(err)
	}

	err = node3.Host.Connect(context.Background(), peer.AddrInfo{node2.Host.ID(), node2.Host.Addrs()})
	if err != nil {
		log.Println(err)
	}

	reply, err := node2.GrpcClient.Greet(context.Background(), node1.Host.ID(), "Dung1")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(reply)

	st1 := node1.Subscribe("thethao", func(st *pubsub.Subscription) {
		for {
			m, e := st.Next(context.Background())
			if e != nil {
				log.Println(e)
			}
			log.Println("st1", string(m.Data))
		}
	})
	st3 := node2.Subscribe("thethao", func(st *pubsub.Subscription) {
		for {
			m, e := st.Next(context.Background())
			if e != nil {
				log.Println(e)
			}
			log.Println("st3", string(m.Data))
		}
	})
	st2 := node3.Subscribe("thethao", func(st *pubsub.Subscription) {
		for {
			m, e := st.Next(context.Background())
			if e != nil {
				log.Println(e)
			}
			log.Println("st2", string(m.Data))
		}
	})

	if st1 != nil {
		panic("Cannot subscribe")
	}
	if st2 != nil {
		panic("Cannot subscribe")
	}
	if st3 != nil {
		panic("Cannot subscribe")
	}

	ticker := time.Tick(1 * time.Second)
	for _ = range ticker {
		log.Println("Publish ...")
		if err = node2.Publish("thethao", []byte("heocon")); err != nil {
			log.Println(err)
		}
	}
}
