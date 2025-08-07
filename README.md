# 🌌 NewCsTeamServer(以下内容为ai生成)

**“在网络的暗****影中，我们重塑规则。”**

欢迎体验 **NewCsTeamServer**，一个用 **Go** 语言重铸的 Cobalt Strike 服务端项目，旨在以极致性能、无缝兼容性和炫酷的浏览器界面统治 C2 战场。这不仅仅是重写——这是 C2 基础设施的革命，专为追求精准、可扩展和风格的操作者打造。
# 构想
```
1.使用golang重构Cs的服务端，并优化相关功能(如分组，数据存储记录，历史记录等)
2.使用web重构ui（缺点:那就是cna的插件化会被遗弃 优点:全部插件和功能dll全放服务端,用户只需要一个浏览器即可）
3.以及更好的用户体验
4.适配原版cs的agent上线与控制
5.读取profile文件，一键适配与启动
....
```
> **📢 提示**：英文版文档请查看 [README_EN.md](README_EN.md)。

---

## 💾 项目愿景

我们致力于打造下一代 TeamServer，超越原版 Cobalt Strike。目标如下：

1. **Go 语言核心**：用 Go 重构服务端，优化分组、数据存储和会话历史记录，性能快如闪电。
2. **Web 界面重塑**：抛弃传统客户端，所有插件和 DLL 部署在服务端，操作者只需一个浏览器即可掌控一切。
   - **优点**：轻量、易访问、无客户端负担。
   - **缺点**：原有的 CNA 插件架构将被弃用，但谁会在意呢？服务端全能掌控！
3. **顶级用户体验**：打造赛博朋克风格的现代、直观界面。
4. **原生 Agent 兼容**：与原版 Cobalt Strike Agent 无缝对接，C2 操作如丝般顺滑。
5. **Profile 解析**：支持 Malleable C2 Profile 文件一键配置与启动。

---

## ⚡ 当前功能

- **Agent 上线**：Agent 与服务端稳定通信，随时待命。
  - 访问 `http://127.0.0.1:8088/get_client_list` 查看客户端列表（简易调试页面）。
- **任务下发**：服务端精准向 Agent 下发任务。
- **结果回传**：Agent 执行任务后将结果回传至服务端。

---

## 🛠 开发进度

项目处于早期开发阶段，但进展迅猛。使用 [geacon](https://github.com/darkr4y/geacon) 进行测试验证，已完成以下里程碑：

### 里程碑
- **Agent 通信**：Agent 与服务端的通信功能已稳定运行。
- **任务管理**：任务下发与结果回传已实现。
- **每日更新**：每天增量优化，持续扩展功能。

### 截图展示
![Agent 通信](png/1.png)  
*Agent 上线，连接稳如磐石。*  
![客户端列表](png/client_list.png)  
*浏览器查看客户端列表，操作即刻上手。*  
![Profile 解析](png/profile.png)  
*Malleable C2 Profile 解析进行中。*  
![Beacon 密钥](png/beacon_key.png)  
*从 `.cobaltstrike.beacon_keys` 提取的 RSA 密钥。*

---

## 📜 更新日志

### 🗓 2025-08-06
- 采用 [Gin](https://github.com/gin-gonic/gin) Web 框架，提升服务端开发效率。
- 实现 Agent 与服务端的核心通信。
- 规划项目文件结构，确保可扩展性和可维护性。
- **后续计划**：
  - 实现 Malleable C2 Profile 解析。
  - 支持加密/解密机制。

### 🗓 2025-08-07
- 集成 [goMalleable](https://github.com/D00Movenok/goMalleable) 实现 Malleable C2 Profile 解析。*从零打造 C2？工作量不小，但我们乐在其中！*
- 使用 [jserial](https://github.com/jkeys089/jserial) 提取 `.cobaltstrike.beacon_keys` 中的 RSA 公私钥，为后续加解密流程铺路。*AI 助力，操作者认证！*

---

## ⏳ 时间线

- **启动时间**：2025-08-05
- **当前阶段**：早期开发，核心通信功能已锁定。
- **未来计划**：Profile 解析、加密支持、完善 Web 界面。

---

## 🕵️‍♂️ 加入行动

我们召集黑客、开发者与梦想家，共同打造 C2 的未来！想加入？以下是你的行动指南：

- **讨论**：访问 [GitHub Issues](https://github.com/your-repo/NewCsTeamServer/issues)，分享你的想法。
- **贡献**：Fork 项目，提交 Pull Request，带来功能增强或修复。
- **反馈**：提出功能建议、用户体验优化，或任何突破常规的创意。

---

## 📜 许可证

本项目采用 [MIT 许可证](LICENSE)。自由 hacking，但请尊重代码。

---

**“低调潜行，迅捷行动，掌控网络。”**  
— NewCsTeamServer 团队