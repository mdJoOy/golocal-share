package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

const (
	broadCastPort     = "9999"
	broadCastAddr     = "255.255.255.255:9999"
	broadCastInterval = 3 * time.Second
)

type PeerInfo struct {
	Id      string `json:"id"`
	TcpPort string `json:"tcpport"`
	Device  string `json:"desktop"`
}

func main() {
	nodeId, err := os.Hostname()
	if err != nil {
		fmt.Println("could not resolve hostname")
		nodeId = "Peer Device"
	}
	self := PeerInfo{Id: nodeId, TcpPort: "8080", Device: "desktop"}

	broadCustPeer(self)
	for {
		listenForPeers(self.Id)

	}
}
func broadCustPeer(self PeerInfo) {
	addr, err := net.ResolveUDPAddr("udp4", broadCastAddr)
	if err != nil {
		log.Fatal("could not resolve broadcast udp address")
	}
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		fmt.Println("couldn't create udp socket")
		return
	}
	ticker := time.NewTicker(broadCastInterval)
	defer ticker.Stop()
	payload, err := json.Marshal(self)
	if err != nil {
		fmt.Println("could not marshal the peer info")
		return
	}
	for range ticker.C {
		_, err := conn.WriteToUDP(payload, addr)
		if err != nil {
			fmt.Println("could not broadcast the peer info", err)
			continue
		}
	}

}
func listenForPeers(selfId string) {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", broadCastPort))
	if err != nil {
		log.Fatal("could not resolve broadcast udp address")
	}
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		fmt.Println("udp listening function problem", err)
	}
	byt := make([]byte, 1024)
	for {
		n, udpAddr, err := conn.ReadFromUDP(byt)
		if err != nil {
			fmt.Println("couldn't read from udp addr")
		}
		var peer PeerInfo
		if err := json.Unmarshal(byt[:n], peer); err != nil {
			fmt.Println("couldn't marshal the peerinfo")
		}
		if peer.Id == selfId {
			continue
		}
		ip := udpAddr.IP.String()
		fmt.Printf("Id: %s, ip: %s, tcp: %, device: %s\n", peer.Id, ip, peer.TcpPort, peer.Device)
	}
}
