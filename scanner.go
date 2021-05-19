package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

func scanPort(protocol, hostname string, port int, wg *sync.WaitGroup) bool {
	defer wg.Done()
	address := hostname + ":" + strconv.Itoa(port)
	conn, err := net.DialTimeout(protocol, address, 500*time.Millisecond)

	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func main() {
	wg := new(sync.WaitGroup)

	args := os.Args[1:]

	lowerPort, upperPort, port := 0, 0, 0
	openPorts := []int{}
	target := args[0]

	lineSeparator := "---------------------"

	if strings.Contains(args[1], "-") {
		// If port range provided, parse out starting and ending port
		portRange := strings.Split(args[1], "-")
		lowerPort, _ = strconv.Atoi(portRange[0])
		upperPort, _ = strconv.Atoi(portRange[1])
	} else {
		port, _ = strconv.Atoi(args[1])
	}

	if !strings.Contains(args[1], "-") {
		// If specific port provided
		wg.Add(1)

		fmt.Printf("[*] Scanning port %d on %s...\n"+lineSeparator+"\n", port, target)

		go func(p int) {
			portStatus := scanPort("tcp", target, p, wg)

			if portStatus == true {
				fmt.Printf("[+]Port %d: open\n", p)
				openPorts = append(openPorts, p)
			} else {
				//fmt.Printf("[-]Port %d: closed\n", p)
			}
		}(port)
	} else {
		// If port range provided
		fmt.Printf("[*] Scanning ports %d-%d on %s...\n"+lineSeparator+"\n", lowerPort, upperPort, target)

		for currentPort := lowerPort; currentPort < upperPort; currentPort++ {
			wg.Add(1)

			go func(p int) {
				portStatus := scanPort("tcp", target, p, wg)

				if portStatus == true {
					fmt.Printf("[+]Port %d: open\n", p)
					openPorts = append(openPorts, p)
				} else {
					//fmt.Printf("[-]Port %d: closed\n", p)
				}
			}(currentPort)

		}

	}

	wg.Wait()

	// Print out all open ports when finished
	fmt.Printf(lineSeparator+"\n[+] Ports open on %s:\n", target)

	for i := 0; i < len(openPorts); i++ {
		fmt.Printf("%d\n", openPorts[i])
	}

}
