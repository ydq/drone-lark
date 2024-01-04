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
		markdown.WriteString("**:SLAP: 失败：** <font color='red'>")
		markdown.WriteString(os.Getenv("DRONE_FAILED_STEPS"))
		markdown.WriteString("</font>\n")
	}

	markdown.WriteString("**:GeneralBusinessTrip: 项目：** [")
	markdown.WriteString(repo)
	markdown.WriteString("](")
	markdown.WriteString(os.Getenv("DRONE_REPO_LINK"))
	markdown.WriteString(")\n")

	if os.Getenv("DRONE_REPO_BRANCH") != "" {
		markdown.WriteString("**:StatusReading: 分支：** <text_tag color='blue'>")
		markdown.WriteString(os.Getenv("DRONE_REPO_BRANCH"))
		markdown.WriteString("</text_tag>\n")
	}

	if os.Getenv("DRONE_TAG") != "" {
		markdown.WriteString("**:Pin: 标签：** <text_tag color='indigo'>")
		markdown.WriteString(os.Getenv("DRONE_TAG"))
		markdown.WriteString("</text_tag>\n")
	}

	author := os.Getenv("DRONE_COMMIT_AUTHOR")
	authorName := os.Getenv("DRONE_COMMIT_AUTHOR_NAME")

	if author == "" {
		if authorName != "" {
			author = authorName
		}
	} else if authorName != "" && author != authorName {
		author = authorName + "@" + author
	}

	if author != "" {
		email := os.Getenv("DRONE_COMMIT_AUTHOR_EMAIL")
		hasEmail := email != ""
		markdown.WriteString("**:EMBARRASSED: 提交：** ")
		if hasEmail {
			markdown.WriteString("[")
		}
		markdown.WriteString(author)
		if hasEmail {
			markdown.WriteString("](mailto:")
			markdown.WriteString(email)
			markdown.WriteString(")")
		}
		markdown.WriteString("\n")
	}

	if os.Getenv("DRONE_COMMIT_SHA") != "" {
		markdown.WriteString("**:Status_PrivateMessage: 信息：** [#")
		markdown.WriteString(os.Getenv("DRONE_COMMIT_SHA")[:8])
		markdown.WriteString("](")
		markdown.WriteString(os.Getenv("DRONE_COMMIT_LINK"))
		markdown.WriteString(")\n")
	}

	markdown.WriteString(" ---\n")
	markdown.WriteString(os.Getenv("DRONE_COMMIT_MESSAGE"))

	elements := []Element{
		{
			Tag:     "markdown",
			Content: markdown.String(),
		},
		{
			Tag: "note",
			Elements: []Element{
				{
					Tag:     "lark_md",
					Content: ":Loudspeaker: [以上信息由 drone 飞书机器人自动发出](" + os.Getenv("DRONE_BUILD_LINK") + ")",
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
	debug := os.Getenv("PLUGIN_DEBUG") == "true"
	if debug {
		fmt.Println("Request Body:", string(jsonBody))
	}

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

	if debug {
		fmt.Println("Response Body:", string(responseBody))
	}
	return nil
}
