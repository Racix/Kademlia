package d7024e

import (
	"log"
	"net"
)

type Network struct {
	me           Contact
	RoutingTable *RoutingTable
}

func (network *Network) SendPingMessage(targetAddress string, rpc *RPCdata) {
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

func (network *Network) Talk(contact *Contact) {
	conn, err := net.DialUDP("udp", *contact.Address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

}

func MarshalRPCdata() {

}
