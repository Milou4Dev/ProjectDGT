package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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
	reset          = "\033[0m"
	bold           = "\033[1m"
	green          = "\033[32m"
	red            = "\033[31m"
)

type Config struct {
	UseConfig    string
	Token        string
	Status       string
	CustomStatus string
	UseEmoji     bool
	EmojiName    string
	EmojiID      string
}

type DiscordUser struct {
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	ID            string `json:"id"`
}

type PresenceUpdate struct {
	Op int `json:"op"`
	D  struct {
		Since      int64      `json:"since"`
		Activities []Activity `json:"activities"`
		Status     string     `json:"status"`
		AFK        bool       `json:"afk"`
	} `json:"d"`
}

type Activity struct {
	Name  string                 `json:"name"`
	Type  int                    `json:"type"`
	State string                 `json:"state,omitempty"`
	Emoji map[string]interface{} `json:"emoji,omitempty"`
}

func main() {
	printHeader()

	fmt.Println("Starting in 10 seconds...")
	time.Sleep(10 * time.Second)

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

	runOnliner(cfg, user, interrupt)
}

func printHeader() {
	fmt.Println(bold + green + "Discord Online Status Setter" + reset)
	fmt.Println("--------------------------------")
	fmt.Println("ProjectDGT by Milou4Dev")
	fmt.Printf("Source code: %s\n\n", projectURL)
}

func getConfig() (Config, error) {
	var cfg Config

	cfg.UseConfig = os.Getenv("USE_CONFIG")
	cfg.Token = os.Getenv("TOKEN")
	cfg.Status = os.Getenv("STATUS")
	cfg.CustomStatus = os.Getenv("CUSTOM_STATUS")
	cfg.UseEmoji = os.Getenv("USE_EMOJI") == "false"
	cfg.EmojiName = os.Getenv("EMOJI_NAME")
	cfg.EmojiID = os.Getenv("EMOJI_ID")

	if strings.ToLower(cfg.UseConfig) == "on" {
		return cfg, nil
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\n" + bold + "Configuration" + reset)

	cfg.Token = promptUntilValid(reader, "Enter your Discord token: ", func(input string) bool { return input != "" })
	cfg.Status = promptUntilValid(reader, "Enter your desired status (online, dnd, idle): ", isValidStatus)
	cfg.CustomStatus, _ = prompt(reader, "Enter your custom status (or press enter for no custom status): ")

	if strings.ToLower(promptUntilValid(reader, "Would you like to use an emoji in your custom status? (y/n): ", func(input string) bool { return input == "y" || input == "n" })) == "y" {
		cfg.UseEmoji = true
		cfg.EmojiName = promptUntilValid(reader, "Enter the emoji name: ", func(input string) bool { return input != "" })
		cfg.EmojiID = promptUntilValid(reader, "Enter the emoji ID: ", func(input string) bool { return input != "" })
	}

	return cfg, nil
}

func promptUntilValid(reader *bufio.Reader, message string, isValid func(string) bool) string {
	attempts := 0
	for {
		input, err := prompt(reader, message)
		if err != nil {
			log.Println(bold + red + "Error:" + reset + " Failed to read input.")
			attempts++
		} else if isValid(input) {
			return input
		} else {
			log.Println(bold + red + "Error:" + reset + " Invalid input.")
			attempts++
		}

		if attempts >= 5 {
			log.Fatalf(bold + red + "Error:" + reset + " Too many invalid attempts. Exiting.")
		}
	}
}

func isValidStatus(status string) bool {
	return status == "online" || status == "dnd" || status == "idle"
}

func prompt(reader *bufio.Reader, message string) (string, error) {
	fmt.Print(message + " ")
	input, err := reader.ReadString('\n')
	return strings.TrimSpace(input), err
}

func fetchUserInfo(token string) (*DiscordUser, error) {
	req, _ := http.NewRequest("GET", discordBaseURL+"/users/@me", nil)
	req.Header.Set("Authorization", token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer closeResponseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token provided, status code: %d", resp.StatusCode)
	}

	var user DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &user, nil
}

func closeResponseBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Printf("Error closing response body: %v", err)
	}
}

func runOnliner(cfg Config, user *DiscordUser, interrupt chan os.Signal) {
	fmt.Printf("\n%sSuccessfully logged in as %s#%s (%s).%s\n", bold+green, user.Username, user.Discriminator, user.ID, reset)
	fmt.Print(bold + "Setting online status... " + reset)

	conn, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		log.Fatalf("Error connecting to Discord gateway: %v", err)
	}
	defer closeConnection(conn)

	heartbeatInterval, err := processHelloMessage(conn)
	if err != nil {
		log.Fatal(err)
	}

	if err := authenticate(conn, cfg.Token); err != nil {
		log.Fatal(err)
	}

	presenceUpdate := createPresenceUpdate(cfg)
	if err := conn.WriteJSON(presenceUpdate); err != nil {
		log.Fatalf("Error sending presence update: %v", err)
	}

	manageHeartbeatAndInterrupt(conn, heartbeatInterval, interrupt, cfg)
}

func closeConnection(conn *websocket.Conn) {
	if err := conn.Close(); err != nil {
		log.Fatalf("Error closing connection: %v", err)
	}
}

func processHelloMessage(conn *websocket.Conn) (time.Duration, error) {
	hello := struct {
		HeartbeatInterval time.Duration `json:"heartbeat_interval"`
	}{}

	if err := conn.ReadJSON(&hello); err != nil {
		return 0, fmt.Errorf("error reading hello message: %v", err)
	}

	if hello.HeartbeatInterval <= 0 {
		return 45 * time.Second, nil
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
			Since      int64      `json:"since"`
			Activities []Activity `json:"activities"`
			Status     string     `json:"status"`
			AFK        bool       `json:"afk"`
		}{
			Since:      time.Now().UnixMilli(),
			Activities: []Activity{},
			Status:     cfg.Status,
			AFK:        false,
		},
	}

	if cfg.CustomStatus != "" {
		activity := Activity{
			Name:  cfg.CustomStatus,
			Type:  4,
			State: cfg.CustomStatus,
		}
		if cfg.UseEmoji {
			activity.Emoji = map[string]interface{}{
				"name": cfg.EmojiName,
				"id":   cfg.EmojiID,
			}
		}
		presenceUpdate.D.Activities = append(presenceUpdate.D.Activities, activity)
	}

	return presenceUpdate
}

func manageHeartbeatAndInterrupt(conn *websocket.Conn, heartbeatInterval time.Duration, interrupt chan os.Signal, cfg Config) {
	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-heartbeatTicker.C:
			if err := conn.WriteJSON(map[string]interface{}{"op": 1, "d": nil}); err != nil {
				log.Printf("Error sending heartbeat: %v", err)
				if err := reconnect(&conn, cfg, heartbeatTicker); err != nil {
					log.Fatalf("Error reconnecting: %v", err)
				}
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

func reconnect(conn **websocket.Conn, cfg Config, heartbeatTicker *time.Ticker) error {
	var err error
	*conn, _, err = websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		return fmt.Errorf("error reconnecting to Discord gateway: %v", err)
	}

	heartbeatInterval, err := processHelloMessage(*conn)
	if err != nil {
		return fmt.Errorf("error processing hello message: %v", err)
	}

	if err := authenticate(*conn, cfg.Token); err != nil {
		return fmt.Errorf("error authenticating: %v", err)
	}

	presenceUpdate := createPresenceUpdate(cfg)
	if err := (*conn).WriteJSON(presenceUpdate); err != nil {
		return fmt.Errorf("error sending presence update: %v", err)
	}

	heartbeatTicker.Reset(heartbeatInterval)
	return nil
}
