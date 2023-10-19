package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	// NewServer
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	network := NewNetwork(*rt)
	if network == nil {
		t.Errorf("NewNetwork return nil")
	}
	var kademlia = NewKademlia(network, 20, 3)
	server := NewServer(kademlia)
	if server == nil {
		t.Error("Server instance nil")
	}
	fmt.Println(kademlia)

	//HandlerRPC
	//
	conn, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0})
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	senderIP := &net.UDPAddr{IP: net.ParseIP("192.168.1.1"), Port: 12345}

	//PING
	SenderID := NewRandomKademliaID()
	var RPCping = NewRPCdata("PING", *SenderID, *SenderID, "", "hej")
	kademlia.HandlerRPC(RPCping, senderIP, conn)
	indexBucket := rt.getBucketIndex(SenderID)
	var contacts = rt.buckets[indexBucket].GetContactAndCalcDistance(SenderID)
	if contacts[0].ID.String() != SenderID.String() {
		t.Error("Contact not correctly added to routingTable")
	}

	// Store
	var RPCstore = NewRPCdata("STORE", *SenderID, *SenderID, "12345", "00112233445566778899AABBCCDDEEFF00112233")
	fmt.Println("crazy")
	kademlia.HandlerRPC(RPCstore, senderIP, conn)
	fmt.Println(kademlia.storeObjects)
	if kademlia.storeObjects[0].data != "00112233445566778899AABBCCDDEEFF00112233" {
		t.Error("String not stored")
	}

	// default
	var defaultRPC = NewRPCdata("hej", *SenderID, *SenderID, "12345", "hej")
	kademlia.HandlerRPC(defaultRPC, senderIP, conn)

	// UnMarshalRPC
	var marshalTest = NewRPCdata("hej", *SenderID, *SenderID, "12345", "hej")
	rpcDataJSON, err := json.Marshal(marshalTest)
	if err != nil {
		t.Error("Error while marshaling: ", err)
	}
	var rpcData1 *RPCdata
	rpcData1, err = UnmarshalRPCdata(rpcDataJSON)
	if err != nil {
		t.Error("Error: ", err)
	}
	if rpcData1.RPCtype != "hej" || rpcData1.SenderID != *SenderID || rpcData1.TargetID != *SenderID {
		t.Error("Wrong data")
	}

	//ContainerServer
	req, err := http.NewRequest("GET", "/whoami", nil)
	if err != nil {
		t.Error()
	}
	_, err = http.NewRequest("Test", "/whoami", nil)
	if err != nil {
		t.Error()
	}
	rr := httptest.NewRecorder()
	http.HandlerFunc(kademlia.MeHandler).ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

	_, err = http.NewRequest("GET", "/get", nil)
	if err != nil {
		t.Error()
	}

	req1, err := http.NewRequest("GET", "/get?hash=asd", nil)
	if err != nil {
		t.Error()
	}
	rr1 := httptest.NewRecorder()
	http.HandlerFunc(kademlia.GetHandler).ServeHTTP(rr1, req1)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, rr.Code)
	}

}
