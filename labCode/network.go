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

func (network *Network) SendPingMessage(targetAddress string, rpc *RPCdata) {
	rpcSend := NewRPCdata("PING", *network.me.ID, rpc.SenderID, rpc.RpcID, "0", nil)
	network.Talk(contact, &rpcSend)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) Pong(contact *Contact, rpc *RPCdata) {
	rpcSend := NewRPCdata("PONG", *network.me.ID, rpc.SenderID, rpc.RpcID, "0", nil)
	network.Talk(contact, &rpcSend)
}
