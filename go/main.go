package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

var logger = log.New(os.Stdout, "go", log.LstdFlags)
var accessToken, accessTokenOK = os.LookupEnv("TRAQ_BOT_ACCESS_TOKEN")

type Payload struct {
	Type  string          `json:"type"`
	ReqID string          `json:"reqId"`
	Body  json.RawMessage `json:"body"`
}

type MessageCreatedBody struct {
	Message struct {
		User struct {
			Name string `json:"name"`
			Bot  bool   `json:"bot"`
		} `json:"user"`
		ChannelID string `json:"channelId"`
		PlainText string `json:"plainText"`
	} `json:"message"`
}

func handleMessage(message []byte) {
	var payload Payload
	err := json.Unmarshal(message, &payload)
	if err != nil {
		panic(err)
	}

	if payload.Type != "MESSAGE_CREATED" {
		logger.Printf("unsupported event(%s): %s\n", payload.ReqID, payload.Type)
		return
	}

	var body MessageCreatedBody
	if err := json.Unmarshal(payload.Body, &body); err != nil {
		logger.Printf("invalid json body(%s): %s\n", payload.ReqID, err.Error())
		return
	}

	if body.Message.User.Bot {
		logger.Printf("bot message(%s)\n", payload.ReqID)
		return
	}

	args := strings.Split(body.Message.PlainText, " ")
	if len(args) != 2 || !strings.HasPrefix(args[0], "@") || args[1] != "go" {
		logger.Printf("invalid args(%s): %s\n", payload.ReqID, body.Message.PlainText)
		return
	}

	stamp := ":golang_new:"
	content := fmt.Sprintf("@%s %s", body.Message.User.Name, stamp)
	if err := postMessage(accessToken, body.Message.ChannelID, content); err != nil {
		logger.Printf("failed to post message(%s): %s\n", payload.ReqID, err.Error())
		return
	}
}

func postMessage(accessToken, channelID, content string) error {
	body, err := json.Marshal(map[string]interface{}{
		"content": content,
		"embed":   true,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://q.trap.jp/api/v3/channels/%s/messages", channelID),
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return nil
}

func main() {
	if !accessTokenOK {
		panic("TRAQ_BOT_ACCESS_TOKEN is not set")
	}

	c, _, err := websocket.DefaultDialer.Dial("wss://q.trap.jp/api/v3/bots/ws", http.Header{
		"Authorization": []string{"Bearer " + accessToken},
	})
	if err != nil {
		panic(err)
	}
	defer c.Close()

	logger.Println("connected")

	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			typ, p, err := c.ReadMessage()
			if err != nil {
				panic(err)
			}

			switch typ {
			case websocket.TextMessage:
				handleMessage(p)
			case websocket.CloseMessage:
				break
			default:
				logger.Printf("unsupported type: %d\n", typ)
			}
		}
	}()

	<-done

	logger.Println("closed")
}
