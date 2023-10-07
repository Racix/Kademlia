package d7024e

import (
	"fmt"
	"strconv"
	"sync"
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
	Network    *Network
	k          int
	alpha      int
	candidates *ContactCandidates
	//responseChan chan []Contact
	storeObjects []StoreObject
}

func NewKademlia(network *Network, k, alpha int) *Kademlia {
	return &Kademlia{
		Network: network,
		k:       k,
		alpha:   alpha,
		candidates: &ContactCandidates{},
		storeObjects: make([]StoreObject, 0),
	}
}

type Candidates struct {
	list	ContactCandidates
	mu sync.Mutex
	newContacts []Contact
}

// type Responses struct {
// 	resChan	chan []Contact
// 	mu sync.Mutex
// }

func (kademlia *Kademlia) LookupContact(target *Contact) ([]Contact, error) {
	candidates := Candidates{}

	visited := make(map[string]bool)

	candidates.mu.Lock()
	candidates.list.Append([]Contact{*target})
	candidates.mu.Unlock()

	// Main lookup loop
	for len(candidates.list.contacts) > 0 {
		fmt.Println("SIZE" + strconv.Itoa(len(candidates.list.contacts)))
		candidates.mu.Lock()
		fmt.Println(candidates.list)
		closest := candidates.list.contacts[0]
		candidates.list.contacts = candidates.list.contacts[1:]
		candidates.mu.Unlock()

		if visited[closest.Address] {
			continue
		}

		visited[closest.Address] = true

		closestContacts := kademlia.Network.routingTable.FindClosestContacts(closest.ID, kademlia.alpha)

		candidates.mu.Lock()
		candidates.list.Append(closestContacts)
		candidates.mu.Unlock()

		responses := make(chan []Contact, kademlia.alpha)

		for i := 0; i < kademlia.alpha; i++ {
			if len(candidates.list.contacts) > 0 {
				candidates.mu.Lock()
				closestContact := candidates.list.contacts[0]
				candidates.list.contacts = candidates.list.contacts[1:]
				candidates.mu.Unlock()

				//for _, c := range closestContact {
				//go kademlia.Network.SendFindContactMessage(&closestContact, responses)
					//responses <- newContacts
				//}


				go func(contact Contact) {
					fmt.Println("THIS STEP 0")
					newContacts := kademlia.Network.SendFindContactMessage(&contact)
					responses <- newContacts
					fmt.Println("THIS STEP 1")
				}(closestContact)
				fmt.Println("THIS STEP 5")
			}
			fmt.Println("THIS STEP 6")
		}
		// check the count!!!
		for i := 0; i < kademlia.alpha; i++ {

			select {
			case resp := <-responses:
				fmt.Println("THIS STEP -1")
				candidates.mu.Lock()
				candidates.newContacts = append(candidates.newContacts, resp...)
				candidates.mu.Unlock()
				fmt.Println("THIS STEP 2")
			default:
				// No more responses expected, exit the loop
				break
			}
		}
		fmt.Println("THIS STEP 3")
		candidates.mu.Lock()
		candidates.list.Append(candidates.newContacts)
		candidates.list.Sort()
		candidates.mu.Unlock()
		fmt.Println("THIS STEP 4")
	}
	fmt.Println("THIS STEP 7")
	// Sort the candidates
	candidates.mu.Lock()
	candidates.list.Sort()
	ck := 0
	if len(candidates.list.contacts) < kademlia.k {ck = len(candidates.list.contacts)} else {ck = kademlia.k}
	kClosestContacts := candidates.list.GetContacts(ck)
	candidates.mu.Unlock()
	fmt.Println("THIS STEP 8")
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
