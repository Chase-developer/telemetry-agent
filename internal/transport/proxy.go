package transport

import (
	"io"
	"log"
	"net"
)

func StartTCPForwarder(listenAddr, forwardAddr string) error {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	log.Printf("Listening on %s, forwarding to %s", listenAddr, forwardAddr)

	for {
		clientConn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go handleConnection(clientConn, forwardAddr)
	}
}

func handleConnection(clientConn net.Conn, forwardAddr string) {
	defer clientConn.Close()

	serverConn, err := net.Dial("tcp", forwardAddr)
	if err != nil {
		log.Printf("Dial error: %v", err)
		return
	}
	defer serverConn.Close()

	// Forward both ways
	go io.Copy(serverConn, clientConn)
	io.Copy(clientConn, serverConn)
}
