package main
 
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)
 
type RequestData struct {
	Message string `json:"message"`
	Ip      string `json:"value"`
}
 
func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading", http.StatusBadRequest)
		return
	}
 
	// Unmarshal the JSON data into a struct
	var requestData RequestData
	if err := json.Unmarshal(body, &requestData); err != nil {
		http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
		return
	}
 
	message := requestData.Message
	ip := requestData.Ip
 
	if ip != "0" {
		sendMessage(message, ip)
		fmt.Fprintf(w, "Forwarding. Message received: "+message+" "+ip)
	} else {
		fmt.Fprintf(w, "Message received: "+message+" "+ip)
		return
	}
}
 
// Function to send a message to an IP address
func sendMessage(message string, url string) ([]byte, error) {
 
	RequestData1 := RequestData{
		Message: message,
		Ip:      "0",
	}
	// Marshal
	packet, err := json.Marshal(RequestData1)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return nil, err
	}
 
	response, err := http.Post("http://"+url+":8080", "application/json", strings.NewReader(string(packet)))
	if err != nil {
		fmt.Println("Error marshaling niinii:", err)
		return nil, err
	}
	defer response.Body.Close()
 
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
 
func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
