package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	discordBaseURL = "https://discord.com/api/v9"
	gatewayURL     = "wss://gateway.discord.gg/?v=9&encoding=json"
	projectURL     = "https://github.com/Milou4Dev/ProjectDGT"
	maxRetries     = 5
	retryDelay     = 5 * time.Second
	timeout        = 10 * time.Second
)

type Config struct {
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
	if err := run(); err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}

func run() error {
	printHeader()
	fmt.Println("Starting in 10 seconds...")
	time.Sleep(10 * time.Second)

	cfg, err := getConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	user, err := fetchUserInfo(cfg.Token)
	if err != nil {
		return fmt.Errorf("failed to fetch user info: %w", err)
	}
	fmt.Println("Success!")

	return runOnliner(cfg, user)
}

func printHeader() {
	fmt.Println("Discord Online Status Setter")
	fmt.Println("--------------------------------")
	fmt.Println("ProjectDGT by Milou4Dev")
	fmt.Printf("Source code: %s\n\n", projectURL)
}

func getConfig() (Config, error) {
	if useConfig := os.Getenv("USE_CONFIG"); strings.ToLower(useConfig) == "on" {
		return Config{
			Token:        os.Getenv("TOKEN"),
			Status:       os.Getenv("STATUS"),
			CustomStatus: os.Getenv("CUSTOM_STATUS"),
			UseEmoji:     os.Getenv("USE_EMOJI") == "true",
			EmojiName:    os.Getenv("EMOJI_NAME"),
			EmojiID:      os.Getenv("EMOJI_ID"),
		}, nil
	}
	return promptForConfig()
}

func promptForConfig() (Config, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("\nConfiguration")

	cfg := Config{}
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
	for attempts := 0; attempts < maxRetries; attempts++ {
		input, err := prompt(reader, message)
		if err != nil {
			log.Println("Error: Failed to read input.")
			continue
		}
		if isValid(input) {
			return input
		}
		log.Println("Error: Invalid input.")
	}
	log.Fatalf("Error: Too many invalid attempts. Exiting.")
	return ""
}

func isValidStatus(status string) bool {
	validStatuses := map[string]bool{"online": true, "dnd": true, "idle": true}
	return validStatuses[status]
}

func prompt(reader *bufio.Reader, message string) (string, error) {
	fmt.Print(message + " ")
	input, err := reader.ReadString('\n')
	return strings.TrimSpace(input), err
}

func fetchUserInfo(token string) (*DiscordUser, error) {
	req, err := http.NewRequest(http.MethodGet, discordBaseURL+"/users/@me", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid token provided, status code: %d", resp.StatusCode)
	}

	var user DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &user, nil
}

func runOnliner(cfg Config, user *DiscordUser) error {
	fmt.Printf("\nSuccessfully logged in as %s#%s (%s).\n", user.Username, user.Discriminator, user.ID)
	fmt.Print("Setting online status... ")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, _, err := websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err != nil {
		return fmt.Errorf("error connecting to Discord gateway: %w", err)
	}
	defer conn.Close()

	heartbeatInterval, err := processHelloMessage(conn)
	if err != nil {
		return err
	}

	if err := authenticate(conn, cfg.Token); err != nil {
		return err
	}

	presenceUpdate := createPresenceUpdate(cfg)
	if err := conn.WriteJSON(presenceUpdate); err != nil {
		return fmt.Errorf("error sending presence update: %w", err)
	}

	return manageHeartbeatAndInterrupt(ctx, conn, heartbeatInterval, cfg)
}

func processHelloMessage(conn *websocket.Conn) (time.Duration, error) {
	var hello struct {
		HeartbeatInterval float64 `json:"heartbeat_interval"`
	}

	if err := conn.ReadJSON(&hello); err != nil {
		return 0, fmt.Errorf("error reading hello message: %w", err)
	}

	if hello.HeartbeatInterval <= 0 {
		return 45 * time.Second, nil
	}

	return time.Duration(hello.HeartbeatInterval) * time.Millisecond, nil
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
		return fmt.Errorf("error sending authentication message: %w", err)
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
			Since:  time.Now().UnixMilli(),
			Status: cfg.Status,
			AFK:    false,
		},
	}

	if cfg.CustomStatus != "" {
		activity := Activity{
			Name:  "Custom Status",
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

func manageHeartbeatAndInterrupt(ctx context.Context, conn *websocket.Conn, heartbeatInterval time.Duration, cfg Config) error {
	heartbeatTicker := time.NewTicker(heartbeatInterval)
	defer heartbeatTicker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-heartbeatTicker.C:
				if err := sendHeartbeat(conn); err != nil {
					if err := reconnect(conn, cfg, heartbeatTicker); err != nil {
						errChan <- fmt.Errorf("error reconnecting: %w", err)
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	select {
	case <-interrupt:
		fmt.Println("\nExiting...")
		cancel := func() {
			ctx.Done()
			wg.Wait()
		}
		defer cancel()
		return conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	case err := <-errChan:
		return err
	}
}

func sendHeartbeat(conn *websocket.Conn) error {
	return conn.WriteJSON(map[string]interface{}{"op": 1, "d": nil})
}

func reconnect(conn *websocket.Conn, cfg Config, heartbeatTicker *time.Ticker) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		conn, _, err = websocket.DefaultDialer.Dial(gatewayURL, nil)
		if err == nil {
			break
		}
		time.Sleep(retryDelay)
	}
	if err != nil {
		return fmt.Errorf("error reconnecting to Discord gateway: %w", err)
	}

	heartbeatInterval, err := processHelloMessage(conn)
	if err != nil {
		return fmt.Errorf("error processing hello message: %w", err)
	}

	if err := authenticate(conn, cfg.Token); err != nil {
		return fmt.Errorf("error authenticating: %w", err)
	}

	presenceUpdate := createPresenceUpdate(cfg)
	if err := conn.WriteJSON(presenceUpdate); err != nil {
		return fmt.Errorf("error sending presence update: %w", err)
	}

	heartbeatTicker.Reset(heartbeatInterval)
	return nil
}
