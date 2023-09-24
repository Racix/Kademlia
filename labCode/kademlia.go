package d7024e

import (
	"fmt"
)

type Kademlia struct {
	network    *Network
	k          int
	alpha      int
	candidates *ContactCandidates
	responseChan chan []Contact // add it here for now
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
			case newContacts := <-kademlia.responseChan:
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

		kademlia.network.SendFindContactMessage(&closest)

		closestContacts := kademlia.network.RoutingTable.FindClosestContacts(target.ID, kademlia.alpha)

		for _, contact := range closestContacts {
			if !alphaContacts[contact.Address] {
				kademlia.candidates.Append([]Contact{contact})
			}
		}
	}
	fmt.Println("Not found")
	return nil
}

// func (kademlia *Kademlia) LookupContact(target *Contact) {
// 	// Initialize a ContactCandidates object to store the closest contacts
// 	candidates := &ContactCandidates{}

// 	// Initialize a set to keep track of visited nodes to avoid loops
// 	visited := make(map[string]bool)

// 	// Initialize a set to keep track of alpha contacts
// 	alphaContacts := make(map[string]bool)

// 	// Initialize a channel to collect responses from Goroutines
// 	responseChan := make(chan []Contact)

// 	// Add the initial contact to the candidates
// 	candidates.Append([]Contact{*target})

// 	processResponses := func() {
// 		var kClosestContacts []Contact
// 		for {
// 			select {
// 			case newContacts := <-kademlia.network.responseChan:
	
// 				// Update the candidates based on the new contacts
// 				kademlia.candidates.Append(newContacts)
	
// 				kademlia.candidates.Sort()
	
// 				// Get the k closest contacts
// 				kClosestContacts = kademlia.candidates.GetContacts(kademlia.k)

// 				fmt.Println("Updated candidates:", kClosestContacts)
// 			}
// 		}
	
// 		return kClosestContacts
// 	}

// 	go processResponses()

// 	// Main lookup loop
// 	for len(candidates.contacts) > 0 {

// 		candidates.Sort()

// 		closest := candidates.contacts[0]

// 		candidates.contacts = candidates.contacts[1:]

// 		if visited[closest.Address] {
// 			continue
// 		}

// 		visited[closest.Address] = true

// 		alphaContacts[closest.Address] = true

// 		kademlia.network.SendPingMessage(&closest)

// 		if closest.ID.Equals(target.ID) {
// 			fmt.Printf("Found target contact: %s\n", closest.String())
// 			return
// 		}

// 		kademlia.network.SendFindContactMessage(&closest)

// 		go func(contact Contact) {
// 			newContacts := kademlia.network.GetContactsFromFindNode(&contact)

// 			responseChan <- newContacts
// 		}(closest)

// 		closestContacts := kademlia.network.routingTable.FindClosestContacts(target.ID, kademlia.alpha)

// 		for _, contact := range closestContacts {
// 			if !alphaContacts[contact.Address] {
// 				candidates.Append([]Contact{contact})
// 			}
// 		}

// 	}

// 	fmt.Println("Target contact not found")
// }

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
	fmt.Println(hash)
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
	fmt.Println(string(data))
}
