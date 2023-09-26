package d7024e

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
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

func (kademlia *Kademlia) HandlerRPC(RPC *RPCdata, senderIP string) {
	switch RPC.RPCtype {
	case "IncPING":
		//NOTE: distance for the new contact?
		theContact := NewContact(&RPC.SenderID, senderIP)
		theContact.Address = senderIP
		indexBucket := kademlia.network.RoutingTable.getBucketIndex(&RPC.SenderID)
		kademlia.network.RoutingTable.buckets[indexBucket].AddContact(theContact)

		kademlia.network.Pong(&theContact, RPC)

	case "OutPing":
		kademlia.network.SendPingMessage(kademlia.network.me.Address, RPC)

	case "IncSTORE":

	case "OutSTORE":

	case "FIND_VALUE":
		//FIND_VALUE

	case "FIND_NODE":
		//FIND_NODE

	default:
		// defualt
	}
	return
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
