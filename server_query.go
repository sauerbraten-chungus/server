package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/sauerbraten/extinfo"
)

type AuthResponse struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

type Player struct {
	Name     string `json:"name"`
	Frags    int    `json:"frags"`
	Deaths   int    `json:"deaths"`
	Accuracy int    `json:"accuracy"`
}

type ServerQueryClient struct {
	httpClient      *http.Client
	server          *extinfo.Server
	playerServiceIP string
	authServiceIP   string
	token           string
	apiKey          string
}

func (sqc *ServerQueryClient) exportMatchData() {
	clients, err := sqc.server.GetAllClientInfo()
	if err != nil {
		fmt.Println(err)
		return
	}

	var players []Player
	for _, client := range clients {
		fmt.Printf("== Client %d ==\n", client.ClientNum)
		fmt.Printf("Ping: %d\n", client.Ping)
		fmt.Printf("Name: %s\n", client.Name)
		fmt.Printf("Team: %s\n", client.Team)
		fmt.Printf("Frags: %d\n", client.Frags)
		fmt.Printf("Deaths: %d\n", client.Deaths)
		fmt.Printf("Team Kills: %d\n", client.Teamkills)
		fmt.Printf("Accuracy: %d\n", client.Accuracy)

		player := Player{
			Name:     client.Name,
			Frags:    client.Frags,
			Deaths:   client.Deaths,
			Accuracy: client.Accuracy,
		}

		players = append(players, player)
	}

	bodyBytes, err := json.Marshal(players)
	if err != nil {
		fmt.Println(err)
		return
	}
	bodyReader := bytes.NewBuffer(bodyBytes)

	url := fmt.Sprintf("%s/players/batch", sqc.playerServiceIP)
	req, err := http.NewRequest("POST", url, bodyReader)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sqc.token)

	fmt.Println(sqc.token)

	resp, err := sqc.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("respBody: %s\n", respBody)
}

func (sqc *ServerQueryClient) obtainJWT() (string, error) {
	url := fmt.Sprintf("%s/auth", sqc.authServiceIP)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("Error creating request: %w", err)
	}
	req.Header.Set("CHUNGUS-KEY", sqc.apiKey)

	resp, err := sqc.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error getting response: %w", err)
	}
	defer resp.Body.Close()

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return "", fmt.Errorf("Error decoding response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Auth error: %s", authResp.Error)
	}

	return authResp.Token, nil
}

func NewServerQueryClient(serverIP, playerServiceIP, authServiceIP, apiKey string, port int) (*ServerQueryClient, error) {
	serverAddr := net.UDPAddr{
		IP:   net.ParseIP(serverIP),
		Port: port,
	}

	server, err := extinfo.NewServer(serverAddr, 5*time.Second)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	sqc := &ServerQueryClient{
		server:          server,
		httpClient:      httpClient,
		playerServiceIP: playerServiceIP,
		authServiceIP:   authServiceIP,
		apiKey:          apiKey,
	}

	token, err := sqc.obtainJWT()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	sqc.token = token

	return sqc, nil
}
