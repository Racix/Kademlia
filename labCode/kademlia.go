package d7024e

import (
	"fmt"
	//"strconv"
	"sync"
	//"errors"
)

type StoreObject struct {
	rpcID            string
	data             string
	dataLength       int
	key              *KademliaID
	senderKademliaID KademliaID
}

func NewStoreObject(rpcID string, data string, dataLength int, key *KademliaID, senderKademliaID KademliaID) *StoreObject {
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
	closestContacts ContactCandidates
	mu sync.Mutex
	newContacts []Contact
	closest Contact
}

// type Responses struct {
// 	resChan	chan []Contact
// 	mu sync.Mutex
// }

func visit(v map[string]bool, c []Contact) bool{
	for _,contact := range c {
		if !v[contact.Address] {
			return true
		}
	}
	return false
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) ([]Contact, error) {
	candidates := Candidates{}

	visited := make(map[string]bool)


	closestContacts := kademlia.Network.routingTable.FindClosestContacts(target, kademlia.k)
	candidates.mu.Lock()
	candidates.closest = Contact{NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"),"0.0.0.0",NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")} 
	//candidates.list.Append([]Contact{*target})
	candidates.list.Append(closestContacts)
	candidates.mu.Unlock()

	// Main lookup loop
	// fmt.Println("SIZE" + strconv.Itoa(len(candidates.list.contacts)))

	responses := make(chan []Contact)
	var closestContact Contact

	for candidates.closestContacts.Len() < kademlia.k {
		if candidates.closestContacts.Len() > 0 {
			if candidates.closestContacts.contacts[0].Less(&candidates.closest) {
				candidates.closest = candidates.closestContacts.contacts[0]

			} else if candidates.list.Len() < 1 {
				// fmt.Println("DEAD HERE!!!")
				break
			}

		}
		var alpha int
		for i := 0; i < kademlia.alpha; i++ {
			if candidates.list.Len() > 0 {
				candidates.mu.Lock()
				// fmt.Printf("BEFORE: %v\n", candidates.list.contacts)
				closestContact = candidates.list.contacts[0]
				candidates.list.contacts = candidates.list.contacts[1:]
				// fmt.Printf("AFTER: %v\n", candidates.list.contacts)
				candidates.mu.Unlock()
				// fmt.Println("AA")
				if visited[closestContact.Address] {
					//fmt.Printf("THE IS FINAL RESULT: (%v,%v)\n",closest.Address,visited[closest.Address] )
					continue
				}
				// fmt.Println("BB")
		
				visited[closestContact.Address] = true
				alpha++
				//for _, c := range closestContact {
				//go kademlia.Network.SendFindContactMessage(&closestContact, responses)
					//responses <- newContacts
				//}

				// go func(contact Contact, responses chan []Contact) {
				// 	fmt.Println("THIS STEP 0")
				// 	newContacts := kademlia.Network.SendFindContactMessage(&contact)
				// 	candidates.mu.Lock()
				// 	contact.CalcDistance(target.ID)
				// 	candidates.closestContacts.Append([]Contact{contact})
				// 	candidates.mu.Unlock()
				// 	responses <- newContacts
				// 	fmt.Println("THIS STEP 1")
				// }(closestContact, responses)


				go func(contact Contact) {
					// fmt.Println("THIS STEP 0")
					newContacts := kademlia.Network.SendFindContactMessage(&contact)
					candidates.mu.Lock()
					contact.CalcDistance(target)
					candidates.closestContacts.Append([]Contact{contact})
					candidates.mu.Unlock()
					responses <- newContacts
					// fmt.Println("THIS STEP 1")
				}(closestContact)
				// fmt.Println("THIS STEP 5")
			}
			// fmt.Println("THIS STEP 6")
		}
		var newContacts []Contact
		for i := 0; i < alpha; i++ /*visit(visited,candidates.list.contacts)*/ {
			//fmt.Println("THIS STEP -1")

			// fmt.Println("THIS STEP -1 " + strconv.Itoa(candidates.list.Len()))
			//candidates.mu.Lock()
			newContacts = append(newContacts, <-responses...)
			//candidates.mu.Unlock()
			// fmt.Println("THIS STEP 2")

			// select {
			// case resp := <-responses:
			// 	fmt.Println("THIS STEP -1 " + strconv.Itoa(candidates.list.Len()))
			// 	candidates.mu.Lock()
			// 	newContacts = append(newContacts, resp...)
			// 	candidates.mu.Unlock()
			// 	fmt.Println("THIS STEP 2")
			// default:
			// 	// No more responses expected, exit the loop
			// 	fmt.Printf("SKIPP!!!: %v\n", closestContact)
			// 	break
			// }

		}

		// fmt.Println("THIS STEP 3")
		candidates.mu.Lock()
		// fmt.Printf("NEWCONTACTS: %v\n", newContacts)
		candidates.list.Append(newContacts)
		// fmt.Printf("NEWclosestContacts: %v\n", candidates.closestContacts.contacts)
		candidates.closestContacts.Sort()
		// for _, c := range candidates.closestContacts {
		// }
		candidates.mu.Unlock()
		// fmt.Println("THIS STEP 4")
	}
	// fmt.Println("THIS STEP 7")
	// Sort the candidates
	candidates.mu.Lock()
	// fmt.Println(candidates.newContacts)
	// closestContacts = kademlia.Network.routingTable.FindClosestContacts(target.ID, kademlia.k)
	// candidates.list.Append(closestContacts)
	// candidates.list.Sort()
	candidates.closestContacts.Sort()
	ck := 0
	if len(candidates.closestContacts.contacts) < kademlia.k {ck = len(candidates.closestContacts.contacts)} else {ck = kademlia.k}
	kClosestContacts := candidates.closestContacts.GetContacts(ck)
	candidates.mu.Unlock()
	// fmt.Println("THIS STEP 8")
	return kClosestContacts, nil
}


func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
	fmt.Println(hash)
}

// func (kademlia *Kademlia) Store(data []byte) {
// 	// TODO
// 	fmt.Println(string(data))
// }

// In i Kademlia
func (kademlia *Kademlia) Store(data string) {
	if len(data) > 255 {
		//errors.New("Value too big")
		fmt.Print("Value too big")
		return
	} else {
		contacts, err := kademlia.LookupContact(NewKademliaID(data))
		if err != nil {
			fmt.Print("dead here: ",err)
		}
		kademlia.Network.SendStoreMessage(data, &contacts[0])
	}
}
