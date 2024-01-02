package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Element struct {
	Tag      string    `json:"tag,omitempty"`
	Content  string    `json:"content,omitempty"`
	Elements []Element `json:"elements,omitempty"`
	URL      string    `json:"url,omitempty"`
	Type     string    `json:"type,omitempty"`
	Text     *Element  `json:"text,omitempty"`
	Actions  []Element `json:"actions,omitempty"`
}

type Header struct {
	Template string  `json:"template,omitempty"`
	Title    Element `json:"title,omitempty"`
}

type Card struct {
	Elements []Element `json:"elements,omitempty"`
	Header   Header    `json:"header,omitempty"`
}

type Body struct {
	Timestamp int64  `json:"timestamp"`
	Sign      string `json:"sign"`
	MsgType   string `json:"msg_type"`
	Card      Card   `json:"card"`
}

func main() {
	webhook := os.Getenv("PLUGIN_WEBHOOK")
	if webhook == "" {
		fmt.Println("缺少 webhook 配置")
		return
	}

	secret := os.Getenv("PLUGIN_SECRET")
	if secret == "" {
		fmt.Println("缺少 secret 配置")
		return
	}

	timestamp := time.Now().Unix()
	sign := generateSignature(timestamp, secret)

	repo := os.Getenv("DRONE_REPO_NAME")

	var title strings.Builder
	var color string

	if os.Getenv("DRONE_BUILD_STATUS") == "success" {
		color = "green"
		title.WriteString("✅ ")
		title.WriteString(repo)
		title.WriteString(" 构建成功 #")
	} else {
		color = "red"
		title.WriteString("❌ ")
		title.WriteString(repo)
		title.WriteString(" 构建失败 #")
	}
	title.WriteString(os.Getenv("DRONE_BUILD_NUMBER"))

	header := Header{
		Template: color,
		Title: Element{
			Tag:     "plain_text",
			Content: title.String(),
		},
	}

	var markdown strings.Builder

	if os.Getenv("DRONE_FAILED_STEPS") != "" {
		markdown.WriteString("**🙅🏻‍♂️ 失败：** ")
		markdown.WriteString(os.Getenv("DRONE_FAILED_STEPS"))
		markdown.WriteString("\n")
	}

	markdown.WriteString("**📦 项目：** [")
	markdown.WriteString(repo)
	markdown.WriteString("](")
	markdown.WriteString(os.Getenv("DRONE_REPO_LINK"))
	markdown.WriteString(")\n")

	if os.Getenv("DRONE_REPO_BRANCH") != "" {
		markdown.WriteString("**🖇️ 分支：** ")
		markdown.WriteString(os.Getenv("DRONE_REPO_BRANCH"))
		markdown.WriteString("\n")
	}

	if os.Getenv("DRONE_TAG") != "" {
		markdown.WriteString("**🏷️ 标签：** ")
		markdown.WriteString(os.Getenv("DRONE_TAG"))
		markdown.WriteString("\n")
	}

	if os.Getenv("DRONE_COMMIT_AUTHOR") != "" {
		markdown.WriteString("**👤 提交：** [")
		hasNick := os.Getenv("DRONE_COMMIT_AUTHOR_NAME") != "" && os.Getenv("DRONE_COMMIT_AUTHOR") != os.Getenv("DRONE_COMMIT_AUTHOR_NAME")
		if hasNick {
			markdown.WriteString(os.Getenv("DRONE_COMMIT_AUTHOR_NAME"))
			markdown.WriteString("@")
		}
		markdown.WriteString(os.Getenv("DRONE_COMMIT_AUTHOR"))
		markdown.WriteString("](mailto:")
		markdown.WriteString(os.Getenv("DRONE_COMMIT_AUTHOR_EMAIL"))
		markdown.WriteString(")\n")
	}

	if os.Getenv("DRONE_COMMIT_SHA") != "" {
		markdown.WriteString("**📝 信息：** [#")
		markdown.WriteString(os.Getenv("DRONE_COMMIT_SHA")[:8])
		markdown.WriteString("](")
		markdown.WriteString(os.Getenv("DRONE_COMMIT_LINK"))
		markdown.WriteString(")\n")
	}

	markdown.WriteString("\n---\n")
	markdown.WriteString(os.Getenv("DRONE_COMMIT_MESSAGE"))

	elements := []Element{
		{
			Tag:     "markdown",
			Content: markdown.String(),
		},
		{
			Tag: "action",
			Actions: []Element{
				{
					Tag:  "button",
					Type: "primary",
					URL:  os.Getenv("DRONE_BUILD_LINK"),
					Text: &Element{
						Tag:     "plain_text",
						Content: "去 Drone 查看本次构建详情",
					},
				},
			},
		},
		// {
		// 	Tag: "hr",
		// },
		{
			Tag: "note",
			Elements: []Element{
				{
					Tag:     "plain_text",
					Content: "🪧 以上信息由 drone 飞书机器人自动发出",
				},
			},
		},
	}

	body := Body{
		Timestamp: timestamp,
		Sign:      sign,
		MsgType:   "interactive",
		Card: Card{
			Header:   header,
			Elements: elements,
		},
	}

	err := sendRequest(webhook, body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func generateSignature(timestamp int64, secret string) string {
	message := fmt.Sprintf("%v\n%v", timestamp, secret)
	mac := hmac.New(sha256.New, []byte(message))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return signature
}

func sendRequest(url string, body Body) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	// fmt.Println(string(jsonBody))

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	// 读取响应主体内容
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println("Response Body:", string(responseBody))
	return nil
}
