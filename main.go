package main

import (
	"fmt"
	"net"
	"time"

	"github.com/sauerbraten/extinfo"
)

func main() {
	serverAddr := net.UDPAddr{
		IP:   net.ParseIP("192.168.1.3"),
		Port: 28785,
	}

	server, err := extinfo.NewServer(serverAddr, 5*time.Second)
	if err != nil {
		fmt.Println(err)
		return
	}

	clients, err := server.GetAllClientInfo()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, client := range clients {
		fmt.Printf("== Client %d ==\n", client.ClientNum)
		fmt.Printf("Ping: %d\n", client.Ping)
		fmt.Printf("Name: %s\n", client.Name)
		fmt.Printf("Team: %s\n", client.Team)
		fmt.Printf("Frags: %d\n", client.Frags)
		fmt.Printf("Deaths: %d\n", client.Deaths)
		fmt.Printf("Team Kills: %d\n", client.Teamkills)
		fmt.Printf("Accuracy: %d\n", client.Accuracy)
	}
}
