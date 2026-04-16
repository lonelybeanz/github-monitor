package monitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func SendNotification(repo string, r *GitHubRelease) {
	webhookURL := os.Getenv("NOTIFY_WEBHOOK")
	if webhookURL == "" {
		return
	}

	// 1. 拼接标题和正文进行检查
	fullText := r.Name + " " + r.Body
	isImportant, hitKeyword := CheckKeywords(fullText)

	// 2. 构造差异化消息
	var title, contentPrefix string

	if isImportant {
		// 🚨 重点提醒样式
		title = fmt.Sprintf("🚨🚨🚨 [%s] 发生重要事件: %s", repo, r.TagName)
		contentPrefix = fmt.Sprintf("⚠️ 触发关键词: **%s**\n\n", strings.ToUpper(hitKeyword))
	} else {
		// 📦 普通更新样式
		title = fmt.Sprintf("📦 [%s] 发布新版本: %s", repo, r.TagName)
		contentPrefix = ""
	}

	// 3. 构造最终 JSON (适配飞书 Lark interactive card)
	// 注意：Release Note 可能很长，建议截断，防止消息发送失败
	shortBody := r.Body
	if len(shortBody) > 500 {
		shortBody = shortBody[:500] + "..."
	}

	// 飞书消息格式
	msg := map[string]interface{}{
		"msg_type": "interactive", // 飞书使用交互式卡片
		"card": map[string]interface{}{
			"header": map[string]interface{}{
				"title": map[string]interface{}{
					"tag":     "plain_text",
					"content": title,
				},
				"template": func() string {
					if isImportant {
						return "red" // 重要消息用红色
					}
					return "blue" // 普通消息用蓝色
				}(),
			},
			"elements": []map[string]interface{}{
				{
					"tag": "div",
					"text": map[string]interface{}{
						"tag": "lark_md",
						"content": fmt.Sprintf("%s**📅 时间:** %s\n\n**📝 说明:**\n%s\n\n[🔗 点击查看详情](%s)",
							contentPrefix,
							r.PublishedAt.Format(time.DateTime),
							shortBody,
							r.HTMLURL),
					},
				},
			},
		},
	}

	payload, _ := json.Marshal(msg)
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("❌ 通知发送失败: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("📨 通知已发送: %s (重要性: %v)", title, isImportant)
}
