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
				kademlia.Store([]byte(strings.Join(get[1:], " ")))
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
		case "send":
			contacts := kademlia.Network.routingTable.FindClosestContacts(kademlia.Network.routingTable.me.ID, 20)
			for i := range contacts {
				fmt.Println(contacts[i].Address + " vs " + get[1])
				if (contacts[i].Address == get[1]) {
					fmt.Println(contacts[i].String())
					//kademlia.LookupContact(&contacts[i])
					kademlia.Network.SendFindContactMessage(&contacts[i])
				}
			}
		case "look":
			contacts := kademlia.Network.routingTable.FindClosestContacts(kademlia.Network.routingTable.me.ID, 20)
			for i := range contacts {
				fmt.Println(contacts[i].Address + " vs " + get[1])
				if (contacts[i].Address == get[1]) {
					fmt.Println(contacts[i].String())
					fmt.Println(kademlia.LookupContact(&contacts[i]))
					//kademlia.Network.SendFindContactMessage(&contacts[i])
				}
			}
		default:
			fmt.Println("Not an option!")
		}
	}

}