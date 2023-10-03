package d7024e

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

type RPCdata struct {
	RPCtype  string     `json:"rpcType"`
	SenderID KademliaID `json:"senderID"`
	TargetID KademliaID `json:"targetID"`
	RpcID    string     `json:"rpcID"`
	Value    string     `json:"value"`
	//for FIND_node & Node lookup
	Contacts []Contact `json:"contacts"`
}

func NewRPCdata(rpcType string, senderID KademliaID, targetID KademliaID, rpcID string, value string) *RPCdata {
	var genID string
	if rpcID == "" {
		genID = GenerateID()
	} else {
		genID = rpcID
	}
	return &RPCdata{
		RPCtype:  rpcType,
		SenderID: senderID,
		TargetID: targetID,
		RpcID:    genID, //here
		Value:    value,
		Contacts: make([]Contact, 0),
	}
}

func GenerateID() string {
	randomBytes := make([]byte, 20)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal(err)
	}
	randomID := hex.EncodeToString(randomBytes)
	return randomID
}
