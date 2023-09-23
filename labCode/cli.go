package d7024e

import (
	"fmt"
	"os"
	"strings"
	"bufio"
)

func (kademlia *Kademlia) Cli() {
	//var i string
	var line string
	var err error
	fmt.Println("Give input")
	for {
		fmt.Print("input>: ")
		c := make(chan int)
		go func() {
			scanner := bufio.NewReader(os.Stdin)
			line, err = scanner.ReadString('\n')
			line = strings.TrimSuffix(line, "\n")
			//fmt.Scan(&i)
			if err != nil {
				fmt.Print(err)
			}
			c <- 0
		}()

		<- c
		// trimmed := strings.TrimSpace(line)
		get := strings.Split(line, " ")

		switch get[0] {
		case "put":
			fmt.Println("make som put!")
			kademlia.Store([]byte(strings.Join(get[1:], " ")))
		case "get":
			fmt.Println("make som get!")
			kademlia.LookupData(strings.Join(get[1:], " "))
		case "exit":
			fmt.Println("make som exit!")
			os.Exit(0)
		default:
			fmt.Println("not an option!")
		}
	}

}