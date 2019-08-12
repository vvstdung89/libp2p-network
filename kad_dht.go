package main

import (
	"context"
	"log"
	"node/p2p"

	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
)

// func ConnectToBootstrapNodes(ctx context.Context, h host.Host, mas []multiaddr.Multiaddr) (numConnected int32) {
// 	var wg sync.WaitGroup
// 	for _, ma := range mas {
// 		wg.Add(1)
// 		go func(ma multiaddr.Multiaddr) {
// 			pi, err := pstore.InfoFromP2pAddr(ma)
// 			if err != nil {
// 				panic(err)
// 			}
// 			defer wg.Done()
// 			err = h.Connect(ctx, *pi)
// 			if err != nil {
// 				log.Printf("error connecting to bootstrap node %q: %v", ma, err)
// 			} else {
// 				atomic.AddInt32(&numConnected, 1)
// 			}
// 		}(ma)
// 	}
// 	wg.Wait()
// 	return
// }

func DHTDemo(hosts []*p2p.Node) {
	if len(hosts) <= 0 {
		return
	}

	listNodeDHT := []*dht.IpfsDHT{
		hosts[0].DHT,
	}

	for i := 1; i < len(hosts); i++ {
		prevNode := hosts[i-1]
		node := hosts[i]
		err := node.Host.Connect(context.Background(), peer.AddrInfo{prevNode.Host.ID(), prevNode.Host.Addrs()})
		if err != nil {
			log.Println("connect node error: ", err)
		}
		dhtObject := node.DHT
		// update DHT data
		dhtObject.BootstrapSelf(context.Background())
		listNodeDHT = append(listNodeDHT, dhtObject)
	}

	// get 2 node to check Ping
	nodeA, nodeB := listNodeDHT[0], listNodeDHT[len(listNodeDHT)-1]
	log.Printf("node A: ID: %s, key: (%x)", nodeA.PeerID().Pretty(), nodeA.PeerKey())
	log.Printf("node B: ID: %s, key: (%x)", nodeB.PeerID().Pretty(), nodeB.PeerKey())
	err := nodeA.Ping(context.Background(), nodeB.PeerID())
	if err != nil {
		log.Println("Ping error: ", err)
	} else {
		log.Printf("Success Connect From %s To %s ", nodeA.PeerID(), nodeB.PeerID())
	}
}