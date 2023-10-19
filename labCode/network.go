package d7024e

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
	"sync"
)

const (
	waitTimeout   = 50 * time.Millisecond //5 * time.Second
	receiveBuffer = 10240
)

type Network struct {
	routingTable RoutingTable
	mu sync.Mutex
}

func NewNetwork(routingTable RoutingTable) *Network {
	return &Network{
		routingTable: routingTable,
	}
}


func (network *Network) Talk(contact *Contact, rpcSend *RPCdata) RPCdata {
	udpAddr, err := net.ResolveUDPAddr("udp", contact.Address)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	rpcDataJSON, err := MarshalRPCdata(rpcSend)
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.Write(rpcDataJSON)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, receiveBuffer)
	conn.SetReadDeadline(time.Now().Add(waitTimeout))
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		//return
		log.Fatal(err)
	}

	respRPC, err := UnmarshalRPCdata(buffer[:n])
	if err != nil {
		fmt.Println(string(buffer[:n]))
		fmt.Printf("Error unmarshaling response: %v\n", err)
		//return
		log.Fatal(err)
	}

	fmt.Printf("THE RESPONE FROM FIND_NODE: %v\n", respRPC.Contacts)

	return *respRPC

}

func MarshalRPCdata(data *RPCdata) ([]byte, error) {
	rpcDataJSON, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return rpcDataJSON, err
}

func (network *Network) SendPingMessage(contact *Contact) {
	rpcSend := NewRPCdata("PING", *network.routingTable.me.ID, *contact.ID, "", "This is a PING")
	network.Talk(contact, rpcSend)
}

// FIND_NODE
func (network *Network) SendFindContactMessage(contact *Contact /*, res chan []Contact*/) []Contact {
	rpcSend := NewRPCdata("FIND_NODE", *network.routingTable.me.ID, *contact.ID, "", "This is a FIND_NODE")
	res := network.Talk(contact, rpcSend).Contacts
	network.mu.Lock()
	network.routingTable.AddContact(*contact)
	network.mu.Unlock()
	return res
}

// FIND_VALUE
func (network *Network) SendFindDataMessage(hash string, contact *Contact) string{
	rpcSend := NewRPCdata("FIND_VALUE", *network.routingTable.me.ID, *contact.ID, "", hash)
	res := network.Talk(contact, rpcSend).Value
	return res
}

// STORE
func (network *Network) SendStoreMessage(data string, contact Contact) {
	rpcSend := NewRPCdata("STORE", *network.routingTable.me.ID, *contact.ID, "", data)
	network.Talk(&contact, rpcSend)
}

//func (network *Network) Pong(contact *Contact, rpc *RPCdata) {
//	rpcSend := NewRPCdata("PONG", *network.routingTable.me.ID, rpc.SenderID, rpc.RpcID, "0")
//	network.Talk(contact, rpcSend)
//}
