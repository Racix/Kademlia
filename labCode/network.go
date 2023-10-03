package d7024e

import (
	"encoding/json"
	"log"
	"net"
)

type Network struct {
	me           Contact
	RoutingTable *RoutingTable
}

func (network *Network) Talk(contact *Contact, rpcSend *RPCdata) {
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

}

func MarshalRPCdata(data *RPCdata) ([]byte, error) {
	rpcDataJSON, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return rpcDataJSON, err
}

func (network *Network) SendPingMessage(contact *Contact) {
	rpcSend := NewRPCdata("PING", *network.me.ID, *contact.ID, "", "This is a PING")
	network.Talk(contact, rpcSend)
}

// FIND_NODE
func (network *Network) SendFindContactMessage(contact *Contact) *[]Contact {
	rpcSend := NewRPCdata("FIND_NODE", *network.me.ID, *contact.ID, "", "This is a FIND_NODE")
	network.Talk(contact, rpcSend)

	return &rpcSend.Contacts
}

// FIND_VALUE
func (network *Network) SendFindDataMessage(hash string) *string {
	rpcSend := NewRPCdata("FIND_VALUE", *network.me.ID, *network.me.ID, "", hash)
	network.Talk(&network.me, rpcSend)

	return &rpcSend.Value
}

// STORE
func (network *Network) SendStoreMessage(data string) {
	rpcSend := NewRPCdata("STORE", *network.me.ID, *network.me.ID, "", data)
	network.Talk(&network.me, rpcSend)
}

func (network *Network) Pong(contact *Contact, rpc *RPCdata) {
	rpcSend := NewRPCdata("PONG", *network.me.ID, rpc.SenderID, rpc.RpcID, "0")
	network.Talk(contact, rpcSend)
}
