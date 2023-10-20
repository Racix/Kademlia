package d7024e

import (
	"testing"
	"time"
)

func TestRademlia(t *testing.T) {

	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8080"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "localhost:8080"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8080"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8080"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8080"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8080"))
	rt.AddContact(NewContact(NewKademliaID("3111111400000000000000000000000000000000"), "localhost:8080"))

	net := NewNetwork(*rt)
	kad := NewKademlia(net,0,3)

	kad.Store("qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq")
	kad.LookupData("q")
	kad.LookupContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"))
	kad.StartLookUp()
	kad = NewKademlia(net,1,0)
	go kad.LookupContact(NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"))
	time.Sleep(3 * time.Second)
	kad.k = 0
	kad.Store("qqqqqqqqqqqqqqq")

	GetLocalIPAddress()

}