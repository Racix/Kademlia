package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
)
type StoreObject struct {
	rpcID            string
	data             string
	dataLength       int
	key              string
	senderKademliaID string
}

func NewStoreObject(rpcID string, data string, key string, senderKademliaID string) *StoreObject {
	storeObject := StoreObject{
		rpcID:            rpcID,
		data:             data,
		dataLength:       len(data),
		key:              key,
		senderKademliaID: senderKademliaID,
	}
	return &storeObject
}
storeObjects: make([]StoreObject, 0),
type Server struct {
	kademlia *Kademlia
}

func NewServer(kademlia *Kademlia) *Server {
	server := Server{
		kademlia: kademlia,
	}
	return &server
}

func (kademlia *Kademlia) HandlerRPC(RPC *RPCdata, senderIP string) {

	switch RPC.RPCtype {
	case "PING":
		//NOTE: distance for the new contact?
		theContact := NewContact(&RPC.SenderID, senderIP)
		theContact.Address = senderIP
		indexBucket := kademlia.network.RoutingTable.getBucketIndex(&RPC.SenderID)
		kademlia.network.RoutingTable.buckets[indexBucket].AddContact(theContact)

		kademlia.network.Pong(&theContact, RPC)
	case "PONG":
		log.Println("PONG recieved")
	case "FIND_VALUE":
		hasher := sha1.New()
		hasher.Write([]byte(RPC.Value))
		theHash := hex.EncodeToString(hasher.Sum(nil))
		for i := 0; i < len(kademlia.storeObjects); i++ {
			if theHash == kademlia.storeObjects[i].key {
				RPC.Value = kademlia.storeObjects[i].data
			}
		}

	// LookUpContact uses FIND_NODE
	case "FIND_NODE":
		closestContacts := kademlia.network.RoutingTable.FindClosestContacts(&RPC.TargetID, 3)
		RPC.Contacts = closestContacts
	case "STORE":
		newStoreObject := kademlia.NewStoreObject()
		kademlia.storeObjects = append(kademlia.storeObjects, newStoreObject)
	default:
		// defualt
	}
}

func (kademlia *Kademlia) Listen(ip string, port int) {
	addr, err := net.ResolveUDPAddr("udp", ip+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, senderAddr, err := conn.ReadFromUDP(buffer)
		senderIP := senderAddr.IP.String()
		if err != nil {
			return
		}
		data := buffer[:n]
		fmt.Printf("Received data from %s: %s\n", addr, data)

		unMarshalledData, err := UnmarshalRPCdata(data)
		if err != nil {
			// for now
			return
		}
		go kademlia.HandlerRPC(unMarshalledData, senderIP)
	}

}

func UnmarshalRPCdata(data []byte) (*RPCdata, error) {
	var rpcData RPCdata
	err := json.Unmarshal(data, &rpcData)
	if err != nil {
		return nil, err
	}
	return &rpcData, nil
}
