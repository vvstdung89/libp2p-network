package main

import (
	"context"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"log"
	"node/p2p"
	"time"
)

func main() {

	var node []*p2p.Node
	for i := 0; i < 100; i++ {
		node = append(node, p2p.NewNode(p2p.NodeConfig{Port: 10000 + i, PublicIP: "127.0.0.1", Version: "1.1", MaxConnection: 100}))
	}
	DHTDemo(node)

	return
	//err := node2.Host.Connect(context.Background(), peer.AddrInfo{node1.Host.ID(), node1.Host.Addrs()})
	//node3.Host.Connect(context.Background(), peer.AddrInfo{node2.Host.ID(), node2.Host.Addrs()})
	//node4.Host.Connect(context.Background(), peer.AddrInfo{node3.Host.ID(), node3.Host.Addrs()})
	//
	//if err != nil {
	//	log.Println(err)
	//}

	//go func() {
	//	ticker := time.Tick(1 * time.Second)
	//	for _ = range ticker {
	//		reply, err := node2.GrpcClient.Greet(context.Background(), node1.Host.ID(), "Dung1")
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		log.Println(reply)
	//	}
	//}()
	//go func() {
	//	ticker := time.Tick(1 * time.Second)
	//	for _ = range ticker {
	//		reply, err := node2.GrpcClient.Greet(context.Background(), node1.Host.ID(), "Dung2")
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		log.Println(reply)
	//	}
	//}()
	//go func() {
	//	ticker := time.Tick(1 * time.Second)
	//	for _ = range ticker {
	//		reply, err := node2.GrpcClient.Greet(context.Background(), node1.Host.ID(), "Dung3")
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		log.Println(reply)
	//	}
	//}()

	//node1.Pubsub.RegisterTopicValidator("thethao", node1.DiscardDuplicateMessageValidator("thethao", 10))
	//node2.Pubsub.RegisterTopicValidator("thethao", node2.DiscardDuplicateMessageValidator("thethao", 10))
	//node3.Pubsub.RegisterTopicValidator("thethao", node3.DiscardDuplicateMessageValidator("thethao", 10))
	//node4.Pubsub.RegisterTopicValidator("thethao", node4.DiscardDuplicateMessageValidator("thethao", 10))
	//
	node[0].Subscribe("thethao", func(st *pubsub.Subscription) {
		for {
			m, e := st.Next(context.Background())
			if e != nil {
				log.Println(e)
			}
			log.Println("st1", string(m.Data))
		}
	})
	node[40].Subscribe("thethao", func(st *pubsub.Subscription) {
		for {
			m, e := st.Next(context.Background())
			if e != nil {
				log.Println(e)
			}
			log.Println("st2", string(m.Data))
		}
	})

	//
	//_ = node3.Subscribe("thethao", func(st *pubsub.Subscription) {
	//	for {
	//		m, e := st.Next(context.Background())
	//		if e != nil {
	//			log.Println(e)
	//		}
	//		log.Println("st3", string(m.Data))
	//	}
	//})
	//
	//_ = node4.Subscribe("thethao", func(st *pubsub.Subscription) {
	//	for {
	//		m, e := st.Next(context.Background())
	//		if e != nil {
	//			log.Println(e)
	//		}
	//		log.Println("st4", string(m.Data))
	//	}
	//})
	//
	ticker := time.Tick(1 * time.Second)
	for _ = range ticker {
		log.Println("Publish ...")
		if err := node[5].Publish("thethao", []byte("heocon")); err != nil {
			log.Println(err)
		}
	}
}
