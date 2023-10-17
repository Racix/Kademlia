package d7024e

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type containerServer struct {
	kademlia *Kademlia
}

func (kademlia *Kademlia) StartServer() {
	http.HandleFunc("/get/hash", kademlia.GetHandler)
	http.HandleFunc("/put/string", kademlia.PutHandler)

	if err := http.ListenAndServe("localhost:8081", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func (kademlia *Kademlia) GetHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	kademlia.LookupData(string(body))
}

func (kademlia *Kademlia) PutHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	kademlia.Store(string(body))
}
