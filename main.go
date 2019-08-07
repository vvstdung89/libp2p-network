package main

import (
	"context"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"node/p2p"
	"time"
)

func main() {
	node1 := p2p.NewNode(p2p.NodeConfig{Port: 10000, PublicIP: "127.0.0.1"})

	node2 := p2p.NewNode(p2p.NodeConfig{Port: 10001, PublicIP: "127.0.0.1"})
	err := node2.Host.Connect(context.Background(), peer.AddrInfo{node1.Host.ID(), node1.Host.Addrs()})
	if err != nil {
		log.Println(err)
	}

	go func() {
		ticker := time.Tick(1 * time.Second)
		for _ = range ticker {
			reply, err := node2.GrpcClient.Greet(context.Background(), node1.Host.ID(), "Dung1")
			if err != nil {
				log.Fatal(err)
			}
			log.Println(reply)
		}
	}()
	go func() {
		ticker := time.Tick(1 * time.Second)
		for _ = range ticker {
			reply, err := node2.GrpcClient.Greet(context.Background(), node1.Host.ID(), "Dung2")
			if err != nil {
				log.Fatal(err)
			}
			log.Println(reply)
		}
	}()
	go func() {
		ticker := time.Tick(1 * time.Second)
		for _ = range ticker {
			reply, err := node2.GrpcClient.Greet(context.Background(), node1.Host.ID(), "Dung3")
			if err != nil {
				log.Fatal(err)
			}
			log.Println(reply)
		}
	}()

	node1.Subscribe("thethao", func(st *pubsub.Subscription) {
		for {
			m, e := st.Next(context.Background())
			if e != nil {
				log.Println(e)
			}
			log.Println("st1", string(m.Data))
		}
	})
	node2.Subscribe("thethao", func(st *pubsub.Subscription) {
		for {
			m, e := st.Next(context.Background())
			if e != nil {
				log.Println(e)
			}
			log.Println("st3", string(m.Data))
		}
	})
	log.Println("Publish ...")
	if err = node2.Publish("thethao", []byte("heocon")); err != nil {
		log.Println(err)
	}
	ticker := time.Tick(1 * time.Second)
	for _ = range ticker {
		log.Println("Publish ...")
		if err = node2.Publish("thethao", []byte("heocon")); err != nil {
			log.Println(err)
		}
	}
}
