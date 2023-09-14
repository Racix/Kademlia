package d7024e

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strconv"
)

type Network struct {
	me           Contact
	RoutingTable *RoutingTable
}

type RPCdata struct {
	RPCtype  string    `json:"rpcType"`
	SenderID string    `json:"senderID"`
	TargetID string    `json:"targetID"`
	Value    string    `json:"value"`
	Contacts []Contact `json:"contacts"`
}

func Listen(ip string, port int) {
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
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			return
		}
		data := buffer[:n]
		fmt.Printf("Received data from %s: %s\n", addr, data)

		unMarshalledData, err := UnmarshalRPCdata(data)
		if err != nil {
			return
		}
		go HandlerRPC(unMarshalledData)
	}

}

func HandlerRPC(RPC *RPCdata) {

	switch RPC.RPCtype {
	case "PING":
		// PING
	case "STORE":
		// STOR

	case "FIND_VALUE":
		//FIND_VALUE

	case "FIND_NODE":
		//FIND_NODE

	default:
		// defualt
	}
	return
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
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

func MarshalRPCdata() {

}

func UnmarshalRPCdata(data []byte) (*RPCdata, error) {
	var rpcData RPCdata
	err := json.Unmarshal(data, &rpcData)
	if err != nil {
		return nil, err
	}
	return &rpcData, nil
}

func (network *Network) Talk(contact *Contact) {
	conn, err := net.DialUDP("udp", contact)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

}
