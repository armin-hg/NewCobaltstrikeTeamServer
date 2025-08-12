```
   ______
  /      \ 
 /________\  NewCsTeamServer
 |  0101  |  C2 Domination
 |  C  S  |  Hack the Network
 |________|
  10101010
```

# 🌌 NewCsTeamServer

**“在网络暗影中，重塑 C2 规则。”**

**NewCsTeamServer**：用 Go 重铸 Cobalt Strike 服务端，统治 C2 战场。无缝 Agent 控制、赛博朋克 Web 界面、Malleable Profile 解析，智慧 hacking，掌控网络。

> **📢 提示**：英文版文档请查看 [README_EN.md](README_EN.md)。

---

## 💾 项目愿景

打造超越原版 Cobalt Strike 的下一代 TeamServer，目标如下：

1. **Go 核心**：重构服务端，优化分组、数据存储和会话历史，性能如电。
2. **Web 界面**：抛弃传统客户端，插件与 DLL 服务端部署，浏览器即控。
    - **优点**：轻量、无客户端负担。
    - **缺点**：CNA 插件退役，服务端全能更胜。
3. **顶级 UX**：赛博朋克风格，直观高效 UI。
4. **Agent 兼容**：无缝对接原版 Cobalt Strike Agent，C2 操作顺滑。
5. **Profile 解析**：一键适配 Malleable C2 Profile。
6. **命令扩展**：支持更多原版命令，提升 C2 灵活性。

---

## ⚡ 当前功能

- **Agent 上线**：稳定通信，随时待命。
    - 访问 `http://127.0.0.1:8088/get_client_list` 查看客户端列表。
- **任务下发**：精准下发任务。
- **结果回传**：Agent 执行后回传结果。
- **流量适配**：支持原版 Agent 上线（Stageless 模式）。

---

## 🛠 开发进度

早期开发，进展迅猛。使用 [geacon](https://github.com/darkr4y/geacon) 测试，已完成以下里程碑：

### 里程碑
- **Agent 通信**：稳定 Agent 与服务端通信。
- **任务管理**：任务下发与结果回传。
- **流量适配**：原版 Agent Stageless 上线。
- **每日更新**：持续优化功能。

### 截图展示
![Agent 通信](png/1.png)  
*Agent 上线，稳如磐石。*  
![客户端列表](png/client_list.png)  
*浏览器客户端列表，即刻上手。*  
![Profile 解析](png/profile.png)  
*Malleable C2 Profile 解析中。*  
![Beacon 密钥](png/beacon_key.png)  
*提取 `.cobaltstrike.beacon_keys` RSA 密钥。*  
![流量适配](png/服务端流量适配.png)  
*适配原版 Agent 上线流量。*  
![监听器](png/监听器.png)  
*监听器配置，精准控制。*

---

## 📜 更新日志

### 🗓 2025-08-06
- 采用 [Gin](https://github.com/gin-gonic/gin) 框架，高效开发。
- 实现 Agent 与服务端通信。
- 规划项目结构，保障扩展性。
- **后续**：Profile 解析、加密支持。

### 🗓 2025-08-07
- 集成 [goMalleable](https://github.com/D00Movenok/goMalleable) 解析 Profile。*C2 重建，挑战巨大，乐在其中！*
- 使用 [jserial](https://github.com/jkeys089/jserial) 提取 `.cobaltstrike.beacon_keys` 公私钥。*AI 助力，操作者认证！*

### 🗓 2025-08-12
- 适配 Profile（如 `jquery-c2.4.5.profile`），感谢 [geacon_pro](https://github.com/your-repo/geacon_pro) 的加解密支持。
- 实现原版 Agent 上线（Stageless，相同 Profile 和密钥）。(可以直接用当前目录下的Beacon.exe(原版cobaltstrike生成的Stageless)进行上线测试，如果你相信的话。)
- 配置监听器，同步上线域名与端口。
- **后续**：
    - 开发赛博朋克风格 Web UI，优化交互。
    - 适配更多原版命令（如 `shell`、`upload` 等）。

---

## ⏳ 时间线

- **启动**：2025-08-05
- **当前**：早期开发，核心通信与流量适配完成。
- **未来**：Web UI 开发、命令扩展、Profile 增强。

---

## 🕵️‍♂️ 加入行动

召集黑客、开发者、梦想家，打造 C2 未来！行动指南：

- **讨论**：[GitHub Issues](https://github.com/your-repo/NewCsTeamServer/issues) 分享想法。
- **贡献**：Fork，提交 Pull Request，增强功能。
- **反馈**：提出 UX 优化或突破性创意。

---

## 📜 许可证

采用 [MIT 许可证](LICENSE)。自由 hacking，尊重代码。

---

**“低调潜行，迅捷行动，掌控网络。”**  
— NewCsTeamServer 团队