package d7024e

type RPCdata struct {
	RPCtype  string     `json:"rpcType"`
	SenderID KademliaID `json:"senderID"`
	TargetID KademliaID `json:"targetID"`
	RpcID    string     `json:"rpcID"`
	Value    string     `json:"value"`
	//for FIND_node & Node lookup
	Contacts []Contact `json:"contacts"`
}

func NewRPCdata(rpcType string, senderID KademliaID, targetID KademliaID, rpcID string, value string, contacts []Contact) RPCdata {
	return RPCdata{
		RPCtype:  rpcType,
		SenderID: senderID,
		TargetID: targetID,
		RpcID:    rpcID,
		Value:    value,
		Contacts: contacts,
	}
}
