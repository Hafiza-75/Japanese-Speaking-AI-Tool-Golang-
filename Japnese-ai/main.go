package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
)

const openaiAPI = "https://api.openai.com/v1/chat/completions"
const apiKey = "YOUR_OPENAI_API_KEY"

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func main() {
	fmt.Println("日本語で入力してください (Type in Japanese):")
	var input string
	fmt.Scanln(&input)

	// Step 1: Send input to OpenAI
	reqBody := ChatRequest{
		Model: "gpt-4o-mini", // lightweight and fast
		Messages: []Message{
			{Role: "user", Content: input},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", openaiAPI, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var chatResp ChatResponse
	json.Unmarshal(body, &chatResp)

	output := chatResp.Choices[0].Message.Content
	fmt.Println("🤖 AI: ", output)

	// Step 2: Convert response to Japanese speech using gTTS
	tmpFile := "output.mp3"
	cmd := exec.Command("gtts-cli", output, "--lang", "ja", "--output", tmpFile)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error generating speech:", err)
		return
	}

	// Step 3: Play audio
	exec.Command("start", tmpFile).Run() // Windows
}
