package main

import (
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"net"
)

//go:embed config.json
var configJSON []byte

type Config struct {
	ConnectionPairs []ConnectionPair `json:"connectionPairs"`
}

type ConnectionPair struct {
	Local  string `json:"local"`
	Remote string `json:"remote"`
}

func loadConfig() (*Config, error) {
	var config Config
	err := json.Unmarshal(configJSON, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func runNetc() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	for _, pair := range config.ConnectionPairs {
		go func(local, remote string) {
			tcpAddr, err := net.ResolveTCPAddr("tcp", local)
			checkError(err)
			listener, err := net.ListenTCP("tcp", tcpAddr)
			checkError(err)
			for {
				conn, err := listener.Accept()
				if err != nil {
					continue
				}
				go handleClient(conn, remote)
			}
		}(pair.Local, pair.Remote)
	}

	select {}
}

func handleClient(conn net.Conn, remote string) {
	defer conn.Close()
	remoteConn, err := net.Dial("tcp", remote)
	if err != nil {
		log.Println(err)
		return
	}
	defer remoteConn.Close()
	go io.Copy(conn, remoteConn)
	io.Copy(remoteConn, conn)
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Fatal error ", err.Error())
	}
}
