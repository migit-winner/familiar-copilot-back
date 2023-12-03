package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	COMPLETION_ENDPOINT       = "https://api.openai.com/v1/chat/completions"
	COMPLETION_SYSTEM_MESSAGE = `途中まで書かれた文がbeforeとafterで提供されます。beforeとafterの間を補完する文章をmiddleに作成してください。
	{
	  "before": "",
	  "middle": "",
	  "after": "",
	}
	の形式で出力してください。beforeとmiddleの間に改行が必要な場合はmiddleの先頭に\nを明記してください。`
)

type OpenAIClient struct {
	httpClient *http.Client
	apiKey     string
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return &OpenAIClient{&http.Client{}, apiKey}
}

type completionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type completionRequestBody struct {
	Model    string              `json:"model"`
	Messages []completionMessage `json:"messages"`
}

func (c *OpenAIClient) GenMiddleText(beforeText, afterText string) (string, error) {
	questionMessage, err := json.Marshal(struct {
		Before string `json:"before"`
		After  string `json:"after"`
	}{Before: beforeText, After: afterText})
	if err != nil {
		return "", err
	}

	requestBody := completionRequestBody{
		Model: "gpt-3.5-turbo",
		Messages: []completionMessage{
			{Role: "system", Content: COMPLETION_SYSTEM_MESSAGE},
			{Role: "user", Content: string(questionMessage)},
		},
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, COMPLETION_ENDPOINT, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Printf("openai api response: %s\n", body)

	var respBody struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		fmt.Println("invalid open ai api response")
		return "", err
	}
	fmt.Printf("openai response: %s\n", respBody.Choices[0].Message.Content)

	var middleText struct {
		Middle string `json:"middle"`
	}
	err = json.Unmarshal([]byte(respBody.Choices[0].Message.Content), &middleText)
	if err != nil {
		fmt.Println("invalid open ai response")
		return "", err
	}

	return middleText.Middle, err
}
