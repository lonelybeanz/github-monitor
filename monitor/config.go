package monitor

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// --- 配置区域 ---

// 需要重点提醒的关键词（全部小写）
var AlertKeywords = []string{
	"hardfork",
	"hard fork",
	"security",
	"vulnerability",
	"critical",
	"cve-", // 包含 CVE 漏洞编号
}

// 你可以在这里添加任意数量的仓库
var TargetRepos = []string{
	"ethereum/go-ethereum",
	"bnb-chain/bsc",
}

const StateFile = "state.json"

// State 用于存储所有仓库的最新 Tag: map["owner/repo"] = "tag"
type State map[string]string

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
