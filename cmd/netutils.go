package cmd

import (
	"bufio"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
)

func server(address string) {
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
//SERVER: either data or closes the connection
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
	server(*addr)
}
*/
