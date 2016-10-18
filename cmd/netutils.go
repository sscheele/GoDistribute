package cmd

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var parts map[string][]string

//only run initParts on the central server
func initParts(fileName string, numChunks int) {
	defString := getOwnIP()
	for i := 0; i < numChunks; i++ {
		parts[fmt.Sprintf("%s.part%d", fileName, i)] = []string{defString}
	}
}

func getOwnIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

//behave like a DNS server, distributing the IPs of computers that have various chunk files
func serveDNS(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleDNS(conn)
	}
}

func handleDNS(conn io.ReadWriteCloser) {
	defer conn.Close()
	buf := make([]byte, 1)
	_, err := conn.Read(buf)
	if err != nil {
		return
	}

	if buf[0] == byte("w") {
		buf := make([]byte, 512)

	} else if buf[0] == byte("r") {
		buf := make([]byte, 512)
		reqLen, err := conn.Read(buf)
		if err != nil {
			return
		}
		ips, ok := parts[string(buf)]
		if !ok {
			conn.WriteString("NULL")
			return
		}
		bytes, err := json.Marshal(ips)
		if err != nil {
			fmt.Println("Error marshalling ip array")
			return
		}
		conn.Write(bytes)
	}

}

func serveFiles(address string) {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleIncoming(conn)
	}
}

//expect the following conversation:
//CLIENT: a number representing the chunk they want
//SERVEFILES: either data or closes the connection
func handleIncoming(conn io.ReadWriteCloser) {
	defer conn.Close()

	buf := make([]byte, 20) //accept, at longest, a 20-digit number
	reqLen, err := conn.Read(buf)
	if err != nil {
		return
	}

}

func sendFile(address string, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	filenameSize := int64(len(filename))

	if err = binary.Write(conn, binary.LittleEndian, filenameSize); err != nil {
		log.Println(err)
		return
	}

	if _, err = io.WriteString(conn, filename); err != nil {
		log.Println(err)
		return
	}

	stat, _ := file.Stat()
	if err = binary.Write(conn, binary.LittleEndian, stat.Size()); err != nil {
		log.Println(err)
		return
	}

	br := bufio.NewReader(file)
	bw := bufio.NewWriter(conn)
	defer bw.Flush()

	if _, err = io.CopyN(bw, br, stat.Size()); err != nil {
		log.Println(err)
		return
	}

	log.Println("File sent.")
}

/*
func main() {
	flag.Parse()
	go func() {
		time.Sleep(time.Second)
		sendFile(*addr, "example.txt")
	}()
	serveFiles(*addr)
}
*/
