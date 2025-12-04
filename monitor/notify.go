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

	// 3. 构造最终 JSON (适配钉钉/飞书/Slack markdown)
	// 注意：Release Note 可能很长，建议截断，防止消息发送失败
	shortBody := r.Body
	if len(shortBody) > 500 {
		shortBody = shortBody[:500] + "..."
	}

	msg := map[string]interface{}{
		"msgtype": "markdown", // 建议改用 markdown 以支持加粗
		"markdown": map[string]string{
			"title": title,
			"text": fmt.Sprintf("### %s\n%s📅 时间: %s\n\n📝 说明:\n%s\n\n[🔗 点击查看详情](%s)",
				title,
				contentPrefix,
				r.PublishedAt.Format(time.DateTime),
				shortBody,
				r.HTMLURL),
		},
		// 钉钉特有：at 所有人
		"at": map[string]interface{}{
			"isAtAll": isImportant, // 只有重要事件才 @所有人
		},
	}

	payload, _ := json.Marshal(msg)
	http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	log.Printf("📨 通知已发送: %s (重要性: %v)", title, isImportant)
}
