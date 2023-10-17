package main

import (
	"kademlia-project/labCode"
	"os"
	"log"

)

const (
	hostIp = "10.0.1.2:8080"
	//hostIp = "localhost:8080"
)

func main() {
	var net d7024e.Network
	var kad d7024e.Kademlia
	id := d7024e.NewKademliaID("9f786a0eef5a4b9e0f7dc37212344491e64ccce8")
	ip, err := d7024e.GetLocalIPAddress()
	if err != nil {
		log.Fatal(err)
	}
	//ip = "localhost:8080"
	me := d7024e.NewContact(id, hostIp)
	if ip == hostIp {
		rt := d7024e.NewRoutingTable(me)
		/*rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8080"))
		rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8080"))
		rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8080"))
		rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8001"))
		rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8001"))
		rt.AddContact(d7024e.NewContact(d7024e.NewKademliaID("3111111400000000000000000000000000000000"), "localhost:8001"))*/
		net = *d7024e.NewNetwork(*rt)
		kad = *d7024e.NewKademlia(&net,20,3)
		go kad.Listen(ip/*,completion*/)
	} else {
		net = d7024e.NetworkJoin(&me)
		kad = *d7024e.NewKademlia(&net,20,3)
		go kad.Listen(ip/*,completion*/)
		kad.StartLookUp()
	}

	//completion := make(chan struct{})
	//go kad.Listen(ip/*,completion*/)
	//kad.Network.SendFindContactMessage(&me)	

	go kad.StartServer()
	kad.Cli(os.Stdin)

	//<-completion
}