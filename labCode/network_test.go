package d7024e

import (
	"testing"
)

func TestNetwork(t *testing.T) {
	//func newNetwork
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))
	network := NewNetwork(*rt)
	if network == nil {
		t.Errorf("NewNetwork return nil")
	}

}
