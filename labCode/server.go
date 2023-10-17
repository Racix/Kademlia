package d7024e

import (
	//"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type Server struct {
	kademlia *Kademlia
}

func NewServer(kademlia *Kademlia) *Server {
	server := Server{
		kademlia: kademlia,
	}
	return &server
}

func (kademlia *Kademlia) HandlerRPC(RPC *RPCdata, senderIP *net.UDPAddr, conn *net.UDPConn) {

	switch RPC.RPCtype {
	case "PING":
		//NOTE: distance for the new contact?
		theContact := NewContact(&RPC.SenderID, senderIP.IP.String())
		theContact.Address = senderIP.IP.String()
		indexBucket := kademlia.Network.routingTable.getBucketIndex(&RPC.SenderID)
		kademlia.Network.routingTable.buckets[indexBucket].AddContact(theContact)

		//PONG response
		rpcDataJSON, err := MarshalRPCdata(RPC)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = conn.WriteToUDP(rpcDataJSON, senderIP)

	case "FIND_VALUE":
		fmt.Println("THIS IS THE RIGHT STOP")
		//hasher := sha1.New()
		//hasher.Write([]byte(RPC.Value))
		//theHash := hex.EncodeToString(hasher.Sum(nil))
		fmt.Println("THIS IS THE SIZE",len(kademlia.storeObjects),kademlia.storeObjects)
		for i := 0; i < len(kademlia.storeObjects); i++ {
			//fmt.Println("THIS IS THE KEY",NewKademliaID(RPC.Value),kademlia.storeObjects[i].key)
			if RPC.Value == kademlia.storeObjects[i].key.String() {
				decoded, _ := hex.DecodeString(kademlia.storeObjects[i].key.String())
				RPC.Value = string(decoded)

			}
		}
		// If value is present --> return. Otherwise --> FIND_NODE
		if RPC.Value == "" {
			theContact := NewContact(&RPC.SenderID, fmt.Sprintf("%s:8080", senderIP.IP.String()))
			kademlia.Network.routingTable.AddContact(theContact)
			closestContacts := kademlia.Network.routingTable.FindClosestContacts(&RPC.TargetID, kademlia.k)
			RPC.Contacts = closestContacts
		}
		// Response
		rpcDataJSON, err := MarshalRPCdata(RPC)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = conn.WriteToUDP(rpcDataJSON, senderIP)

	// LookUpContact uses FIND_NODE
	case "FIND_NODE":
		theContact := NewContact(&RPC.SenderID, fmt.Sprintf("%s:8080", senderIP.IP.String()))
		kademlia.Network.routingTable.AddContact(theContact)
		closestContacts := kademlia.Network.routingTable.FindClosestContacts(&RPC.TargetID, kademlia.k)
		RPC.Contacts = closestContacts

		rpcDataJSON, err := MarshalRPCdata(RPC)
		if err != nil {
			log.Fatal(err)
		}

		_, _ = conn.WriteToUDP(rpcDataJSON, senderIP)
	case "STORE":
		newStoreObject := NewStoreObject(RPC.RpcID, RPC.Value, len(RPC.Value), NewKademliaID(RPC.Value), RPC.SenderID)
		kademlia.storeObjects = append(kademlia.storeObjects, *newStoreObject)
		fmt.Println("THIS IS THE SIZE",len(kademlia.storeObjects),kademlia.storeObjects)

		rpcDataJSON, err := MarshalRPCdata(RPC)
		if err != nil {
			log.Fatal(err)
		}
		_, _ = conn.WriteToUDP(rpcDataJSON, senderIP)
	default:
		fmt.Printf("Not a correct RPC type")
	}
}

func (kademlia *Kademlia) Listen(ip string /*,completion chan struct{}*/) {
	//defer close(completion)
	addr, err := net.ResolveUDPAddr("udp", ip)
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
		//senderIP := senderAddr.IP.String()
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
		go kademlia.HandlerRPC(unMarshalledData, senderAddr, conn)
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
