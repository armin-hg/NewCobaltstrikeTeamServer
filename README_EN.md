# 🌌 NewCsTeamServer

**"In the shadows of the network, we rewrite the rules."**

Welcome to **NewCsTeamServer**, a reimagined Cobalt Strike TeamServer forged in **Go**, designed to dominate with enhanced performance, seamless compatibility, and a sleek, browser-based interface. This is not just a rewrite—it's a revolution in C2 infrastructure, crafted for operators who demand precision, scalability, and style.

---

## 💾 Project Vision

We’re building a next-gen TeamServer that outsmarts the original. Our mission:

1. **Go-Powered Core**: Rebuild the Cobalt Strike server in **Go** for blazing-fast performance, optimized grouping, robust data storage, and comprehensive session history tracking.
2. **Web-Based Interface**: Ditch the old client. All plugins and DLLs live server-side—operators only need a browser to control the chaos.
   - **Pros**: Lightweight, accessible, no client-side bloat.
   - **Cons**: Say goodbye to CNA plugins, but who needs them when you’ve got full server-side control?
3. **Elite UX**: A modern, intuitive interface that feels like a cyberpunk dashboard.
4. **Native Agent Compatibility**: Seamless integration with original Cobalt Strike Agents for flawless C2 operations.
5. **Profile Parsing**: One-click configuration with full support for Malleable C2 Profile files.

---

## ⚡ Current Features

- **Agent Uplink**: Agents connect and report in, ready for action.
  - Access the client list at `http://127.0.0.1:8088/get_client_list` for a quick debug view.
- **Task Dispatch**: Issue commands to your Agents with precision.
- **Result Relay**: Agents execute tasks and beam results back to the server.

---

## 🛠 Development Status

We’re in the early stages, but moving fast. Powered by [geacon](https://github.com/darkr4y/geacon) for testing, we’re hammering out core functionality and laying the groundwork for a robust C2 framework.

### Milestones
- **Agent Comms**: Fully functional Agent-to-Server communication.
- **Task Management**: Task dispatch and result retrieval are live.
- **Daily Progress**: Incremental updates, pushing the limits every day.

### Screenshots
![Agent Comms](png/1.png)  
*Agent uplink in action.*  
![Client List](png/client_list.png)  
*Browser-based client list for quick ops.*  
![Profile Parsing](png/profile.png)  
*Malleable C2 Profile parsing in progress.*  
![Beacon Keys](png/beacon_key.png)  
*Extracted RSA keys from `.cobaltstrike.beacon_keys`.*

---

## 📜 Changelog

### 🗓 2025-08-06
- Switched to [Gin](https://github.com/gin-gonic/gin) for high-performance web routing.
- Established core Agent-to-Server communication.
- Structured project for scalability and maintainability.
- **Next Steps**:
  - Parse Malleable C2 Profiles.
  - Implement encryption/decryption workflows.

### 🗓 2025-08-07
- Integrated [goMalleable](https://github.com/D00Movenok/goMalleable) for Malleable C2 Profile parsing. *Building a C2 from scratch? Yeah, it’s a grind, but we’re in it.*
- Used [jserial](https://github.com/jkeys089/jserial) to extract RSA public/private keys from `.cobaltstrike.beacon_keys` for secure Agent comms. *AI-assisted, operator-approved.*

---

## ⏳ Timeline

- **Kickoff**: 2025-08-05
- **Current Phase**: Early development, with core comms locked in.
- **Future**: Profile parsing, encryption, and a polished web UI.

---

## 🕵️‍♂️ Join the Ops

We’re recruiting hackers, coders, and visionaries to shape the future of C2. Got ideas? Skills? Bugs to squash? Here’s how to dive in:

- **Discuss**: Drop into [GitHub Issues](https://github.com/your-repo/NewCsTeamServer/issues) to share your thoughts.
- **Contribute**: Fork, code, and submit Pull Requests with your enhancements.
- **Feedback**: Suggest features, UX improvements, or new ways to break the mold.

---

## 📜 License

This project is licensed under the [MIT License](LICENSE). Hack freely, but respect the code.

---

**"Stay low, move fast, control the network."**  
— NewCsTeamServer Crew