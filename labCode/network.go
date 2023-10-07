package d7024e

import (
	"encoding/json"
	"log"
	"net"
	"time"
	"fmt"
)

const (
	waitTimeout   = 5 * time.Second
	receiveBuffer = 1024
)

type Network struct {
	routingTable RoutingTable
}

func NewNetwork(routingTable RoutingTable) *Network {
	return &Network{
		routingTable: routingTable,
	}
}

func NetworkJoin(contact *Contact) Network {
	id := NewRandomKademliaID()
	ip, err := GetLocalIPAddress()
	if err != nil {
		log.Fatal(err)
	}
	me := NewContact(id, ip)
	rt := NewRoutingTable(me)
	rt.AddContact(*contact)
	return *NewNetwork(*rt)

}

func GetLocalIPAddress() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				return fmt.Sprintf("%s:8080", ipnet.IP.String()), nil
			}
		}
	}

	return "", fmt.Errorf("no suitable IP address found")
}

func (network *Network) Talk(contact *Contact, rpcSend *RPCdata) []Contact{
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
		//fmt.Printf("Error reading response: %v\n", err)
		//return
		log.Fatal(err)
	}

	repJSON, err := UnmarshalRPCdata(buffer[:n])
	if err != nil {
		//fmt.Printf("Error unmarshaling response: %v\n", err)
		//return
		log.Fatal(err)
	}

	fmt.Printf("THE RESPONE FROM FIND_NODE: %v\n", repJSON.Contacts)

	return repJSON.Contacts

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
func (network *Network) SendFindContactMessage(contact *Contact/*, res chan []Contact*/) []Contact{
	rpcSend := NewRPCdata("FIND_NODE", *network.routingTable.me.ID, *contact.ID, "", "This is a FIND_NODE")
	//res <- network.Talk(contact, rpcSend)
	return network.Talk(contact, rpcSend)
}

// FIND_VALUE
func (network *Network) SendFindDataMessage(hash string) *string {
	rpcSend := NewRPCdata("FIND_VALUE", *network.routingTable.me.ID, *network.routingTable.me.ID, "", hash)
	network.Talk(&network.routingTable.me, rpcSend)

	return &rpcSend.Value
}

// STORE
func (network *Network) SendStoreMessage(data string) {
	rpcSend := NewRPCdata("STORE", *network.routingTable.me.ID, *network.routingTable.me.ID, "", data)
	network.Talk(&network.routingTable.me, rpcSend)
}

func (network *Network) Pong(contact *Contact, rpc *RPCdata) {
	rpcSend := NewRPCdata("PONG", *network.routingTable.me.ID, rpc.SenderID, rpc.RpcID, "0")
	network.Talk(contact, rpcSend)
}
