package main

import (
	"encoding/binary"
	"github.com/Hermes-Bird/faraway-test-task.git/internal/challenge"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("server address should be provided")
	}

	log.Printf("trying to dial with %s\n", os.Args[1])
	con, err := net.Dial("tcp", os.Args[1])
	if err != nil {
		log.Fatalf("error while dialing with %s\n", os.Args[1])
	}

	defer con.Close()

	bs := make([]byte, 20)
	log.Println("reading challenge from connection")
	n, err := con.Read(bs)
	if err != nil && err != io.EOF {
		return
	}

	log.Println("starting to solve given challenge")
	res := challenge.SolveChallenge(bs[:n])
	log.Printf("challenge solved (%d), sending to solution to connection\n", binary.BigEndian.Uint32(res))
	_, err = con.Write(res)
	if err != nil {
		log.Println("error while writing solution to connection")
		return
	}
	log.Println("solution sent awaiting result")

	bs = make([]byte, 256)
	builder := strings.Builder{}
	for {
		n, err = con.Read(bs)
		if err != nil && err != io.EOF {
			log.Println("error while reading result from the connection")
			return
		}

		builder.Write(bs[:n])

		if err == io.EOF || n == 0 {
			break
		}
	}

	log.Printf("result of work - %s\n", builder.String())
}
