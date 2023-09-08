package main
 
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"bytes"
)
 
type RequestData struct {
	Message string `json:"message"`
	Ip      string `json:"value"`
}
 

func handler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
        return
    }
 
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Error reading request body", http.StatusBadRequest)
        return
    }
 
    defer r.Body.Close()
 
    var data map[string]interface{}
    if err := json.Unmarshal(body, &data); err != nil {
        http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
        return
    }
 
    receivedMessage := data["message"].(string)
    receivedIp := data["value"].(string)
    sendMessage(receivedMessage, receivedIp)
    if receivedIp == "0" {
    	fmt.Printf("Received message: %s\n", receivedMessage)
    	return
    } else {
    	fmt.Println("Message is sent")
    }

    response := map[string]string{"message": "the message is received"}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
        return
    }
 
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write(jsonResponse)
}
 
func sendMessage(payload string, url string) ([]byte, error) {
    sendData := RequestData{
    	Message: payload,
    	Ip: "0",
    }
    
    client := &http.Client{}
    requestURL := "http://" + url + ":8080"
    sendData2, _ := json.Marshal(sendData)
    requestBody := bytes.NewBuffer([]byte(sendData2))
 
    req, err := http.NewRequest("POST", requestURL, requestBody)
    if err != nil {
        return nil, err
    }
 
    req.Header.Set("Content-Type", "application/json")
 
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
 
    response, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
 
    return response, nil
}
 
func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
