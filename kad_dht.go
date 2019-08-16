package main

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/peer"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"log"
	"node/p2p"
	"time"
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

// CONNECT NODE TO ANOTHER NODE WITH INFO AND UPDATE DHT
func ConnectAndUpdateDHT(node *p2p.Node, addrInfo peer.AddrInfo) (*p2p.Node, error) {
	err := node.Host.Connect(context.Background(), addrInfo)
	if err != nil {
		log.Println("connect node error: ", err)
		return nil, err
	}
	dht := node.DHT
	// update DHT
	dht.BootstrapSelf(context.Background())
	return node, nil
}

// DEMO NODE CONNECT AND CHECK CONNECTION WITH ANOTHER NODE IN NETWORK
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
		dhtObject.BootstrapRandom(context.Background())
		//dhtObject.BootstrapRandom(context.Background())
		//dhtObject.Update(context.Background(), node.Host.ID())
		listNodeDHT = append(listNodeDHT, dhtObject)
	}
	time.Sleep(5 * time.Second)
	hosts[0].DHT.RoutingTable().Print()
	fmt.Println(hosts[0].DHT.RoutingTable().Size())
	for i, v := range hosts[0].Host.Network().Conns() {
		fmt.Println(i+1, v)
	}

	// get 2 node to check Ping
	nodeA, nodeB := listNodeDHT[0], listNodeDHT[len(listNodeDHT)-1]
	log.Printf("node A: ID: %s, key: (%x)", nodeA.PeerID().Pretty(), nodeA.PeerKey())
	log.Printf("node B: ID: %s, key: (%x)", nodeB.PeerID().Pretty(), nodeB.PeerKey())
	//err := nodeA.Ping(context.Background(), nodeB.PeerID())

	//nodeA.Process()

	// PUT VALUE TO NODE
	err := nodeB.PutValue(context.Background(), "/pk/multihash", []byte("123"))
	if err != nil {
		log.Println("put value error: ", err)
	} else {
		log.Println("put value successfully")
	}

	time.Sleep(1 * time.Second)
	log.Println("node A list peers :", nodeA.RoutingTable().ListPeers())
	log.Println("node A peer store:", nodeA.Host().Peerstore().Peers())
	//addr, err := nodeA.FindPeer(context.Background(), nodeB.PeerID())
	value, _ := nodeA.GetValue(context.Background(), "/pk/multihash")
	// fmt.Println(nodeA.GetValue(context.Background(), "a"))
	log.Println("Get value: ", string(value))

}
