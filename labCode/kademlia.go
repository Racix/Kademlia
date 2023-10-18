package d7024e

import (
	//"crypto/sha1"
	"encoding/hex"
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

func (kademlia *Kademlia)StartLookUp() {
	kademlia.LookupContact(kademlia.Network.routingTable.me.ID)
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
	candidates.list.Append(closestContacts)
	candidates.mu.Unlock()

	responses := make(chan []Contact)
	var closestContact Contact

	for candidates.closestContacts.Len() < kademlia.k {
		if candidates.closestContacts.Len() > 0 {
			if candidates.closestContacts.contacts[0].Less(&candidates.closest) {
				candidates.closest = candidates.closestContacts.contacts[0]

			} else if candidates.list.Len() < 1 {
				break
			}

		}
		var alpha int
		for i := 0; i < kademlia.alpha; i++ {
			if candidates.list.Len() > 0 {

				candidates.mu.Lock()
				closestContact = candidates.list.contacts[0]
				candidates.list.contacts = candidates.list.contacts[1:]
				candidates.mu.Unlock()

				if visited[closestContact.Address] {
					continue
				}
		
				visited[closestContact.Address] = true
				alpha++

				go func(contact Contact) {
					newContacts := kademlia.Network.SendFindContactMessage(&contact)
					candidates.mu.Lock()
					contact.CalcDistance(target)
					candidates.closestContacts.Append([]Contact{contact})
					candidates.mu.Unlock()
					responses <- newContacts
				}(closestContact)
			}
		}
		var newContacts []Contact
		for i := 0; i < alpha; i++ {
			newContacts = append(newContacts, <-responses...)
		}
		candidates.mu.Lock()
		candidates.list.Append(newContacts)
		candidates.closestContacts.Sort()
		candidates.mu.Unlock()
	}
	candidates.mu.Lock()
	candidates.closestContacts.Sort()
	ck := 0
	if len(candidates.closestContacts.contacts) < kademlia.k {ck = len(candidates.closestContacts.contacts)} else {ck = kademlia.k}
	kClosestContacts := candidates.closestContacts.GetContacts(ck)
	candidates.mu.Unlock()
	return kClosestContacts, nil
}


// func (kademlia *Kademlia) LookupData(hash string) {
// 	// TODO
// 	fmt.Println(hash)
// }
func (kademlia *Kademlia) LookupData(hash string) string{
	var data string
    if len(hash) != 40 {
        fmt.Println("Too long/short hash")
    } else {
        contacts, err := kademlia.LookupContact(NewKademliaID(hash))
		if err != nil {
			fmt.Print("dead here: ",err)
		}
		fmt.Println("WTF IS HERE")
        data = kademlia.Network.SendFindDataMessage(hash, &contacts[0])
    }
	fmt.Println("IS IT POSSIBLE")
	return data
}

// func (kademlia *Kademlia) Store(data []byte) {
// 	// TODO
// 	fmt.Println(string(data))
// }

// In i Kademlia
func (kademlia *Kademlia) Store(data string) KademliaID{
	var id KademliaID
	if len(data) > 255 {
		//err = errors.New("Value too big")
		fmt.Println("Value too big")
	} else {
		hashValue := hex.EncodeToString([]byte(data)[0:IDLength])
		//hashValue := hex.EncodeToString(sha1.New().Sum([]byte(data))[0:IDLength])
		id = *NewKademliaID(hashValue)
		contacts, err := kademlia.LookupContact(&id)
		if err != nil {
			fmt.Print("dead here: ",err)
		}
		go func() {
			for _, contact := range(contacts) {
				go kademlia.Network.SendStoreMessage(hashValue, contact)
			}
		}()

	}
	return id
}
