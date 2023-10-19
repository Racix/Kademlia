package d7024e

import (
	"fmt"
	"io/ioutil"
	"net/http"
)


func (kademlia *Kademlia) StartServer() {
	http.HandleFunc("/get", kademlia.GetHandler)
	http.HandleFunc("/put/string", kademlia.PutHandler)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func (kademlia *Kademlia) GetHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
        http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
        return
    }

	hash := r.URL.Query().Get("hash")

	// Check if the value parameter is empty
	if hash == "" {
		http.Error(w, "Value parameter is required", http.StatusBadRequest)
		return
	}

	// Handle the value (e.g., display it)
	fmt.Printf("Received value: %s\n", hash)

	data := kademlia.LookupData(string(hash))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Received value: %s\n", data)))

	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	http.Error(w, "Error reading request body", http.StatusInternalServerError)
	// 	return
	// }
	// //kademlia.LookupData(string(body))
	// fmt.Printf("This is get: %v\n", string(body))
}

func (kademlia *Kademlia) PutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
        http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
        return
    }

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	hash := kademlia.Store(string(body))
	//fmt.Printf("This is hash: %v\n", hash)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("This is hash: %s\n", hash.String())))
}
