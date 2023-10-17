package d7024e

import (
	"fmt"
	"os"
	"io"
	"strings"
	"bufio"
)

func (kademlia *Kademlia) Cli(stdin io.Reader) {
	var line string
	var err error
	scanner := bufio.NewReader(stdin)
	fmt.Println("Give input")
	for {
		fmt.Print("input>: ")
		line, err = scanner.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		if err != nil {
			fmt.Print("dead here: ",err)
		}

		// trimmed := strings.TrimSpace(line)
		get := strings.Split(line, " ")

		switch get[0] {
		case "put":
			if(len(get[1:]) > 0) {
				fmt.Println("Make som put!")
				kademlia.Store(strings.Join(get[1:], " "))
			} else {
				fmt.Println("Put takes at LEAST ONE arg")
			}
		case "get":
			if (len(get[1:]) == 1) {
				fmt.Println("Make som get!")
				kademlia.LookupData(get[1])
			} else {
				fmt.Println("Get takes ONE arg")
			}
		case "exit":
			fmt.Println("Make som exit!")
			os.Exit(1)
		case "sendme":
			contacts := kademlia.Network.routingTable.buckets
			fmt.Println(contacts)
		case "send":
			contacts := kademlia.Network.routingTable.FindClosestContacts(kademlia.Network.routingTable.me.ID, 20)
			fmt.Println(contacts)
			for i := range contacts {
				fmt.Println(contacts[i].Address + " vs " + get[1])
				if (contacts[i].Address == get[1]) {
					fmt.Println(contacts[i].String())
					//kademlia.LookupContact(&contacts[i])
					kademlia.Network.SendFindContactMessage(&contacts[i])
				}
			}
		case "look":
			c := NewKademliaID(get[1])//NewContact(NewKademliaID(get[1]),get[2])
			v,_ :=kademlia.LookupContact(c)
			fmt.Printf("THE IS FINAL RESULT: %v\n",v )
		default:
			fmt.Println("Not an option!")
		}
	}

}