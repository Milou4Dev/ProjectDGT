package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	discordBaseURL = "https://discord.com/api/v9"
	gatewayURL     = "wss://gateway.discord.gg/?v=9&encoding=json"
	projectURL     = "https://github.com/Milou4Dev/ProjectDGT"

	reset = "\033[0m"
	bold  = "\033[1m"
	green = "\033[32m"
	red   = "\033[31m"
)

type Config struct {
	Token        string
	Status       string
	CustomStatus string
	EmojiName    string
	EmojiID      string
	UseEmoji     bool
}

type DiscordUser struct {
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	ID            string `json:"id"`
}

type PresenceUpdate struct {
	Op int `json:"op"`
	D  struct {
		Since      int64         `json:"since"`
		Activities []interface{} `json:"activities"`
		Status     string        `json:"status"`
		AFK        bool          `json:"afk"`
	} `json:"d"`
}

func main() {
	printIntro()

	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	user, err := fetchUserInfo(cfg.Token)
	if err != nil {
		log.Fatalf("\n%sFailed%s to fetch user info: %v", bold+red, reset, err)
	}
	fmt.Println(bold + green + "Success!" + reset)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go runOnliner(cfg, user, interrupt)
	<-interrupt
	fmt.Println("\nExiting...")
}

func printIntro() {
	fmt.Println(bold + green + "Discord Online Status Setter" + reset)
	fmt.Println("--------------------------------")
	fmt.Println("ProjectDGT by Milou4Dev")
	fmt.Printf("Source code: %s\n\n", projectURL)
}

func getConfig() (Config, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n" + bold + "Configuration" + reset)

	cfg := Config{}
	var err error

	cfg.Token, err = prompt(reader, "Enter your Discord token: ")
	if err != nil {
		return cfg, err
	}
	cfg.Status, err = prompt(reader, "Enter your desired status (online, dnd, idle): ")
	if err != nil {
		return cfg, err
	}
	cfg.CustomStatus, err = prompt(reader, "Enter your custom status (or type 'none' for no custom status): ")
	if err != nil {
		return cfg, err
	}

	if cfg.Token == "" {
		return cfg, fmt.Errorf(bold + red + "Error:" + reset + " A valid token is required")
	}

	input, err := prompt(reader, "Would you like to use an emoji in your custom status? (y/n): ")
	if err != nil {
		return cfg, err
	}
	if strings.ToLower(input) == "y" {
		cfg.UseEmoji = true
		cfg.EmojiName, err = prompt(reader, "Enter the emoji name: ")
		if err != nil {
			return cfg, err
		}
		cfg.EmojiID, err = prompt(reader, "Enter the emoji ID: ")
		if err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}

func prompt(reader *bufio.Reader, message string) (string, error) {
	fmt.Print(message + " ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func fetchUserInfo(token string) (*DiscordUser, error) {
	req, err := http.NewRequest("GET", discordBaseURL+"/users/@me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token provided")
	}
	user := &DiscordUser{}
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}

func runOnliner(cfg Config, user *DiscordUser, interrupt chan os.Signal) {
	fmt.Printf("\n%sSuccessfully logged in as %s#%s (%s).%s\n", bold+green, user.Username, user.Discriminator, user.ID, reset)

	fmt.Print(bold + "Setting online status... " + reset)
	conn, err := establishWebSocketConnection(gatewayURL)
	if err != nil {
		log.Fatalf("Error connecting to Discord gateway: %v", err)
	}
	defer conn.Close()

	heartbeatInterval, err := processHelloMessage(conn)
	if err != nil {
		log.Fatal(err)
	}

	err = authenticate(conn, cfg.Token)
	if err != nil {
		log.Fatal(err)
	}

	presenceUpdate := createPresenceUpdate(cfg)
	if err := conn.WriteJSON(presenceUpdate); err != nil {
		log.Fatalf("Error sending presence update: %v", err)
	}

	manageHeartbeatAndInterrupt(conn, heartbeatInterval, interrupt)
}

func establishWebSocketConnection(url string) (*websocket.Conn, error) {
	dialer := *websocket.DefaultDialer
	dialer.HandshakeTimeout = 10 * time.Second
	conn, _, err := dialer.Dial(url, nil)
	return conn, err
}

func processHelloMessage(conn *websocket.Conn) (time.Duration, error) {
	hello := struct {
		HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	}{}

	if err := conn.ReadJSON(&hello); err != nil {
		return 0, fmt.Errorf("error reading hello message: %v", err)
	}

	if hello.HeartbeatInterval <= 0 {
		hello.HeartbeatInterval = 45 * time.Second
	}

	return hello.HeartbeatInterval, nil
}

func authenticate(conn *websocket.Conn, token string) error {
	auth := map[string]interface{}{
		"op": 2,
		"d": map[string]interface{}{
			"token": token,
			"properties": map[string]string{
				"$os":      "Windows",
				"$browser": "Chrome",
				"$device":  "ProjectDGT",
			},
			"intents": 513,
		},
	}

	if err := conn.WriteJSON(auth); err != nil {
		return fmt.Errorf("error sending authentication message: %v", err)
	}
	return nil
}

func createPresenceUpdate(cfg Config) PresenceUpdate {
	presenceUpdate := PresenceUpdate{
		Op: 3,
		D: struct {
			Since      int64         `json:"since"`
			Activities []interface{} `json:"activities"`
			Status     string        `json:"status"`
			AFK        bool          `json:"afk"`
		}{
			Since:      time.Now().Unix(),
			Activities: []interface{}{},
			Status:     cfg.Status,
			AFK:        false,
		},
	}

	if cfg.CustomStatus != "none" {
		activity := map[string]interface{}{
			"name":  "Custom Status",
			"type":  4,
			"state": cfg.CustomStatus,
		}
		if cfg.UseEmoji {
			activity["emoji"] = map[string]interface{}{
				"name": cfg.EmojiName,
				"id":   cfg.EmojiID,
			}
		}
		presenceUpdate.D.Activities = append(presenceUpdate.D.Activities, activity)
	}

	return presenceUpdate
}

func manageHeartbeatAndInterrupt(conn *websocket.Conn, heartbeatInterval time.Duration, interrupt chan os.Signal) {
	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-heartbeatTicker.C:
			if err := conn.WriteJSON(map[string]interface{}{"op": 1, "d": nil}); err != nil {
				log.Fatalf("Error sending heartbeat: %v", err)
			}
		case <-interrupt:
			fmt.Println("\nExiting...")
			if err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")); err != nil {
				log.Println("Error closing WebSocket connection:", err)
			}
			return
		}
	}
}
