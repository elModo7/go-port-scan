package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

var (
	author    string
	version   string
	startPort int
	endPort   int
	targetIP  string
)

func banner() {
	name := fmt.Sprintf("go-port-scan (v.%s)", version)
	banner := `
  ________      __________              __   _________                     
 /  _____/  ____\______   \____________/  |_/   _____/ ____ _____    ____  
/   \  ___ /  _ \|     ___/  _ \_  __ \   __\_____  \_/ ___\\__  \  /    \ 
\    \_\  (  <_> )    |  (  <_> )  | \/|  | /        \  \___ / __ \|   |  \
 \______  /\____/|____|   \____/|__|   |__|/_______  /\___  >____  /___|  /
        \/                                         \/     \/     \/     \/ 
	`

	allLines := strings.Split(banner, "\n")
	w := len(allLines[1])

	fmt.Println(banner)
	color.Green(fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(name))/2, name)))
	color.Blue(fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(author))/2, author)))
	fmt.Println()
}

func init() {
	version = "1.0.1"
	author = "lismore, elModo7"
	banner()
}

func main() {
	flag.StringVar(&targetIP, "t", "127.0.0.1", "Target IP")
	flag.IntVar(&startPort, "sp", 20, "Start Port")
	flag.IntVar(&endPort, "ep", 1024, "End Port")
	flag.Parse()

	if startPort < 1 || endPort > 65535 || startPort > endPort {
		log.Fatal("invalid port range. Please specify startPort and endPort between 1 and 65535, with startPort less than or equal to endPort.")
	}

	var wg sync.WaitGroup

	// Limit concurrent connections
	const maxConcurrent = 200
	sem := make(chan struct{}, maxConcurrent)

	for port := startPort; port <= endPort; port++ {
		wg.Add(1)
		go checkPortOpen(&wg, sem, port)
	}

	wg.Wait()
}

func checkPortOpen(wg *sync.WaitGroup, sem chan struct{}, port int) {
	defer wg.Done()

	sem <- struct{}{}
	defer func() { <-sem }()

	address := net.JoinHostPort(targetIP, strconv.Itoa(port))

	conn, err := net.DialTimeout("tcp", address, 800*time.Millisecond)
	if err != nil {
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}(conn)

	log.Printf("Port %d open", port)
}
