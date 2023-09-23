package d7024e

import (
	"fmt"
)

type Kademlia struct {
	network    *Network
	k          int
	alpha      int
	candidates *ContactCandidates
}

func NewKademlia(network *Network, k, alpha int) *Kademlia {
	return &Kademlia{
		network:    network,
		k:          k,
		alpha:      alpha,
		candidates: &ContactCandidates{},
	}
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {

	alphaContacts := make(map[string]bool)
	kademlia.candidates.Append([]Contact{*target})

	processResponses := func() {
		for {
			select {
			case newContacts := <-kademlia.network.responseChan:
				kademlia.candidates.Append(newContacts)
				kademlia.candidates.Sort()

				kClosest := kademlia.candidates.GetContacts(kademlia.k)
				fmt.Println("Updated candidates:", kClosest)
			}
		}
	}

	go processResponses()

	// Main loop
	for len(kademlia.candidates.contacts) > 0 {

		kademlia.candidates.Sort()
		closest := kademlia.candidates.contacts[0]
		kademlia.candidates.contacts = kademlia.candidates.contacts[1:]
		alphaContacts[closest.Address] = true

		kademlia.network.SendPingMessage(&closest)

		// Check if contact found
		if closest.ID.Equals(target.ID) {
			fmt.Printf("Lookup contact: %s\n", closest.String())
			return []Contact{closest}
		}

		kademlia.network.SendFindNodeMessage(&closest)

		closestContacts := kademlia.network.routingTable.FindClosestContacts(target.ID, kademlia.alpha)

		for _, contact := range closestContacts {
			if !alphaContacts[contact.Address] {
				kademlia.candidates.Append([]Contact{contact})
			}
		}
	}
	fmt.Println("Not found")
	return nil
}


func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
