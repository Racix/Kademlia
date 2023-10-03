package d7024e

import (
	"fmt"
)

type StoreObject struct {
	rpcID            string
	data             string
	dataLength       int
	key              string
	senderKademliaID string
}

func NewStoreObject(rpcID string, data string, key string, senderKademliaID string) *StoreObject {
	storeObject := StoreObject{
		rpcID:            rpcID,
		data:             data,
		dataLength:       len(data),
		key:              key,
		senderKademliaID: senderKademliaID,
	}
	return &storeObject
}

type Kademlia struct {
	network    *Network
	k          int
	alpha      int
	candidates *ContactCandidates
	//responseChan chan []Contact
	storeObjects []StoreObject
}

func NewKademlia(network *Network, k, alpha int) *Kademlia {
	return &Kademlia{
		network: network,
		k:       k,
		alpha:   alpha,
		//candidates: &ContactCandidates{},
		storeObjects: make([]StoreObject, 0),
	}
}

func (kademlia *Kademlia) LookupContact(target *Contact) ([]Contact, error) {
	candidates := &ContactCandidates{}

	visited := make(map[string]bool)

	candidates.Append([]Contact{*target})

	// Main lookup loop
	for len(candidates.contacts) > 0 {
		candidates.Sort()
		closest := candidates.contacts[0]
		candidates.contacts = candidates.contacts[1:]

		if visited[closest.Address] {
			continue
		}

		visited[closest.Address] = true

		closestContacts := kademlia.network.RoutingTable.FindClosestContacts(closest.ID, kademlia.alpha)

		candidates.Append(closestContacts)

		responses := make(chan []Contact, kademlia.alpha)

		for i := 0; i < kademlia.alpha; i++ {
			if len(candidates.contacts) > 0 {
				closestContact := candidates.contacts[0]
				candidates.contacts = candidates.contacts[1:]
				go func(contact Contact) {
					//newContacts := kademlia.network.SendFindContactMessage(&contact)
					//responses <- newContacts
				}(closestContact)
			}
		}

		var newContacts []Contact
		for i := 0; i < kademlia.alpha; i++ {
			newContacts = append(newContacts, <-responses...)
		}

		candidates.Append(newContacts)
	}

	// Sort the candidates
	candidates.Sort()

	kClosestContacts := candidates.GetContacts(kademlia.k)

	return kClosestContacts, nil
}

// func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {

// 	alphaContacts := make(map[string]bool)
// 	kademlia.candidates.Append([]Contact{*target})

// 	processResponses := func() {
// 		for {
// 			select {
// 			case newContacts := <-kademlia.responseChan:
// 				kademlia.candidates.Append(newContacts)
// 				kademlia.candidates.Sort()

// 				kClosest := kademlia.candidates.GetContacts(kademlia.k)
// 				fmt.Println("Updated candidates:", kClosest)
// 			}
// 		}
// 	}

// 	go processResponses()

// 	// Main loop
// 	for len(kademlia.candidates.contacts) > 0 {

// 		kademlia.candidates.Sort()
// 		closest := kademlia.candidates.contacts[0]
// 		kademlia.candidates.contacts = kademlia.candidates.contacts[1:]
// 		alphaContacts[closest.Address] = true

// 		kademlia.network.SendPingMessage(&closest)

// 		// Check if contact found
// 		if closest.ID.Equals(target.ID) {
// 			fmt.Printf("Lookup contact: %s\n", closest.String())
// 			return []Contact{closest}
// 		}

// 		kademlia.network.SendFindContactMessage(&closest)

// 		closestContacts := kademlia.network.RoutingTable.FindClosestContacts(target.ID, kademlia.alpha)

// 		for _, contact := range closestContacts {
// 			if !alphaContacts[contact.Address] {
// 				kademlia.candidates.Append([]Contact{contact})
// 			}
// 		}
// 	}
// 	fmt.Println("Not found")
// 	return nil
// }

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
	fmt.Println(hash)
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
	fmt.Println(string(data))
}
