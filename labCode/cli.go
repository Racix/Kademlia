package d7024e

import (
	"fmt"
	"os"
	"strings"
	"bufio"
)

func (kademlia *Kademlia) Cli() {
	var line string
	var err error
	fmt.Println("Give input")
	for {
		fmt.Print("input>: ")

		scanner := bufio.NewReader(os.Stdin)
		line, err = scanner.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		if err != nil {
			fmt.Print(err)
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
			os.Exit(0)
		default:
			fmt.Println("Not an option!")
		}
	}

}