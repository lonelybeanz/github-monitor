package monitor

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// --- 配置区域 ---

// 需要重点提醒的关键词（全部小写）
var AlertKeywords = getAlertKeywordsFromEnv()

// 你可以在这里添加任意数量的仓库
var TargetRepos = getTargetReposFromEnv()

const StateFile = "state.json"

// State 用于存储所有仓库的最新 Tag: map["owner/repo"] = "tag"
type State map[string]string

// 从环境变量获取告警关键词列表，如果未设置则使用默认值
func getAlertKeywordsFromEnv() []string {
	envValue := os.Getenv("ALERT_KEYWORDS")
	if envValue == "" {
		return []string{
			"hardfork",
			"hard fork",
			"security",
			"vulnerability",
			"critical",
			"cve-", // 包含 CVE 漏洞编号
		}
	}
	
	// 按逗号分割环境变量
	keywords := strings.Split(envValue, ",")
	// 清理每个关键词前后的空格
	for i, kw := range keywords {
		keywords[i] = strings.TrimSpace(kw)
	}
	return keywords
}

// 从环境变量获取目标仓库列表，如果未设置则使用默认值
func getTargetReposFromEnv() []string {
	envValue := os.Getenv("TARGET_REPOS")
	if envValue == "" {
		return []string{
			"ethereum/go-ethereum",
			"bnb-chain/bsc",
			"base/node",
			"anza-xyz/agave",
		}
	}
	
	// 按逗号分割环境变量
	repos := strings.Split(envValue, ",")
	// 清理每个仓库名称前后的空格
	for i, repo := range repos {
		repos[i] = strings.TrimSpace(repo)
	}
	return repos
}

func LoadState() State {
	s := make(State)
	data, err := os.ReadFile(StateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return s // 文件不存在，返回空 map
		}
		log.Printf("⚠️ 读取状态文件失败: %v", err)
		return s
	}
	_ = json.Unmarshal(data, &s)
	return s
}

func SaveState(s State) {
	data, _ := json.MarshalIndent(s, "", "  ") // 美化输出
	err := os.WriteFile(StateFile, data, 0644)
	if err != nil {
		log.Printf("❌ 保存状态失败: %v", err)
	} else {
		log.Println("💾 状态文件已更新")
	}
}

// --- 辅助函数：关键词检测 ---
// checkKeywords 返回是否包含关键词，以及具体包含哪一个（用于显示）
func CheckKeywords(text string) (bool, string) {
	lowerText := strings.ToLower(text)
	for _, kw := range AlertKeywords {
		if strings.Contains(lowerText, kw) {
			return true, kw
		}
	}
	return false, ""
}
