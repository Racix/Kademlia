package d7024e

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)


func (kademlia *Kademlia) StartServer() {
	http.HandleFunc("/get", kademlia.GetHandler)
	http.HandleFunc("/whoami", kademlia.MeHandler)
	http.HandleFunc("/put/string", kademlia.PutHandler)

	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func (kademlia *Kademlia) MeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
        http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
        return
    }

	exit := r.URL.Query().Get("")

	data := kademlia.Network.routingTable.me

	if exit == "exit" {
		defer os.Exit(1)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Shout down: %s, %s\n", data.ID.String(), data.Address)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("I am: %s, %s\n", data.ID.String(), data.Address)))

}

func (kademlia *Kademlia) GetHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
        http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
        return
    }

	hash := r.URL.Query().Get("hash")

	if hash == "" {
		http.Error(w, "Value parameter is required", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received value: %s\n", hash)

	data := kademlia.LookupData(string(hash))

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Received value: %s\n", data)))
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("This is hash: %s\n", hash.String())))
}
