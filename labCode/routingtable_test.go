package d7024e

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	//his part added to cover kademliaID more
	id := NewKademliaID("9f786a0eef5a4b9e0f7dc37212344491e64ccce8")
	fmt.Println(NewRandomKademliaID())
	fmt.Println(id)
	id.Less(id)
	//id.Equals(id)

	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("3111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("3111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
		fmt.Println(contacts[i].ID.Less(id))
	}
	// this part is to cover the routingtable more
	contacts2 := rt.FindClosestContacts(id, 5)	
	for i := range contacts2 {
		fmt.Println(contacts[i].String())
	}

	// this part to cover the bucket more
	rt.AddContact(NewContact(NewKademliaID("3111111400000000000000000000000000000000"), "localhost:8002"))
	fmt.Println(rt.buckets[4].Len())
	rt.AddContact(NewContact(id, "localhost:8002"))
}
