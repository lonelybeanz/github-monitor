# 使用 GitHub Actions 监控GitHub中某些仓库的Release更新

**使用 GitHub Actions 维护监控脚本是性价比最高的方案：完全免费、无需维护服务器、自带运行日志**
GitHub Actions 的一个核心特性：无状态（Stateless）。 每次 Action 运行都是在一个全新的、干净的容器里。这意味着脚本运行完，容器销毁，你生成的 state.json 也会消失。
**解决方案**： 我们利用 Git 本身！让 Action 在检测到新版本后，自动将新的 Tag 提交（Commit）并推送（Push）回你的仓库。这样仓库里的文件就充当了“数据库”。


## 第一步：配置通知 Webhook (Secrets)

1. 进入你的 GitHub 仓库页面。
2. 点击顶部菜单 Settings -> 左侧栏 Secrets and variables -> Actions。
3. 点击 New repository secret。
    - Name: NOTIFY_WEBHOOK
    - - Secret: 填入你的机器人 Webhook 地址（例如钉钉机器人的 https://oapi.dingtalk.com/robot/send?access_token=xxxx）。
    - Name: GITHUB_TOKEN（用于避免 API 限流）
    - - Secret: 
4. 保存。

## 第二步：在 GitHub 设置变量

1. 进入仓库 Settings -> Secrets and variables -> Actions -> Variables (注意是 Variables 标签页，不是 Secrets)。
2. 点击 New repository variable。
    - Name: TARGET_REPOS
    - - Value: ethereum/go-ethereum,bnb-chain/bsc,base/node
    - Name: ALERT_KEYWORDS
    - - Value: hardfork,security,critical

## 第三步：验证运行
1. 进入仓库的 Actions 标签页。
2. 在左侧点击 Monitor Geth Releases。
3. 点击右侧的 Run workflow 按钮手动触发一次。
4. 观察结果：
    - 点击运行详情，查看 Log。
    - 如果 state.json 里的版本旧了，Go 程序会识别到新版本，发送通知。
    - Git Commit 步骤会将新的 Tag 写回仓库。
    - 你可以刷新代码页面，看看 state.json 的内容是否变了。