package main

import (
	"fmt"
	"github-monitor/monitor"
	"log"
	"sync"
)

// --- 配置区域 ---

func main() {
	// 1. 加载状态
	state := monitor.LoadState()

	// 2. 并发检查
	var wg sync.WaitGroup
	var mu sync.Mutex // 保护 state 的并发写入
	hasUpdates := false

	log.Printf("🚀 开始检查 %d 个仓库...", len(monitor.TargetRepos))

	for _, repo := range monitor.TargetRepos {
		wg.Add(1)
		go func(repoName string) {
			defer wg.Done()

			// 获取最新 Release
			release, err := monitor.FetchLatestRelease(repoName)
			if err != nil {
				log.Printf("❌ [%s] 获取失败: %v", repoName, err)
				return
			}

			// 检查是否更新
			mu.Lock()
			lastTag := state[repoName]
			mu.Unlock()

			if release.TagName != lastTag {
				fmt.Printf("🎉 [%s] 发现新版本: %s (旧: %s)\n", repoName, release.TagName, lastTag)

				// 发送通知
				monitor.SendNotification(repoName, release)

				// 更新状态 (加锁)
				mu.Lock()
				state[repoName] = release.TagName
				hasUpdates = true
				mu.Unlock()
			} else {
				log.Printf("✅ [%s] 无更新 (%s)", repoName, lastTag)
			}
		}(repo)
	}

	// 等待所有 Goroutine 完成
	wg.Wait()

	// 3. 如果有更新，保存状态文件
	if hasUpdates {
		monitor.SaveState(state)
	} else {
		log.Println("💤 所有仓库均无更新，无需保存状态。")
	}
}
