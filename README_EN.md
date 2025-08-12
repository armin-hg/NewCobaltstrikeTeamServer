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

**"In the shadows of the network, we rewrite C2 rules."**

**NewCsTeamServer**: Go-powered Cobalt Strike server rewrite. Dominate C2 with seamless Agent control, cyberpunk web UI, and Malleable Profile parsing. Hack smarter, rule the network.

---

## 💾 Project Vision

Building a next-gen TeamServer to outsmart Cobalt Strike. Our mission:

1. **Go Core**: Rebuild server for blazing-fast performance, optimized grouping, data storage, and session history.
2. **Web Interface**: Ditch clients. Plugins and DLLs server-side—control via browser.
    - **Pros**: Lightweight, no client bloat.
    - **Cons**: CNA plugins retired, server-side rules supreme.
3. **Elite UX**: Cyberpunk-style, intuitive UI.
4. **Agent Compatibility**: Seamless original Cobalt Strike Agent integration.
5. **Profile Parsing**: One-click Malleable C2 Profile support.
6. **Command Expansion**: Support more native commands for flexible C2.

---

## ⚡ Current Features

- **Agent Uplink**: Reliable Agent-to-Server comms.
    - Access `http://127.0.0.1:8088/get_client_list` for debug client list.
- **Task Dispatch**: Precise command issuance.
- **Result Relay**: Agents return task results.
- **Traffic Compatibility**: Original Agent uplinks (Stageless mode).

---

## 🛠 Development Status

Early development, moving fast. Tested with [geacon](https://github.com/darkr4y/geacon), we’ve hit these milestones:

### Milestones
- **Agent Comms**: Stable Agent-Server communication.
- **Task Management**: Task dispatch and result retrieval.
- **Traffic Compatibility**: Original Agent Stageless uplinks.
- **Daily Progress**: Continuous feature expansion.

### Screenshots
![Agent Comms](png/1.png)  
*Agent uplink in action.*  
![Client List](png/client_list.png)  
*Browser-based client list for quick ops.*  
![Profile Parsing](png/profile.png)  
*Malleable C2 Profile parsing in progress.*  
![Beacon Keys](png/beacon_key.png)  
*Extracted RSA keys from `.cobaltstrike.beacon_keys`.*  
![Traffic Compatibility](png/服务端流量适配.png)  
*Adapted for original Agent uplinks.*  
![Listener](png/监听器.png)  
*Listener config for precise control.*

---

## 📜 Changelog

### 🗓 2025-08-06
- Adopted [Gin](https://github.com/gin-gonic/gin) for high-performance routing.
- Established Agent-Server comms.
- Structured project for scalability.
- **Next**: Profile parsing, encryption.

### 🗓 2025-08-07
- Integrated [goMalleable](https://github.com/D00Movenok/goMalleable) for Profile parsing. *Building C2 from scratch? A grind, but we’re in.*
- Used [jserial](https://github.com/jkeys089/jserial) to extract RSA keys from `.cobaltstrike.beacon_keys`. *AI-assisted, operator-approved.*

### 🗓 2025-08-12
- Adapted Profiles (e.g., `jquery-c2.4.5.profile`), leveraging [geacon_pro](https://github.com/your-repo/geacon_pro)’s encryption.
- Enabled original Agent uplinks (Stageless, same Profile/keys).
- Configured listeners for domain/port alignment.
- **Next**:
    - Develop cyberpunk web UI for seamless interaction.
    - Expand command support (e.g., `shell`, `upload`).

---

## ⏳ Timeline

- **Kickoff**: 2025-08-05
- **Current**: Early development, core comms and traffic compatibility locked.
- **Future**: Web UI, command expansion, Profile enhancements.

---

## 🕵️‍♂️ Join the Ops

Recruiting hackers, coders, visionaries to shape C2’s future. Your mission:

- **Discuss**: Hit [GitHub Issues](https://github.com/your-repo/NewCsTeamServer/issues) with ideas.
- **Contribute**: Fork, submit Pull Requests with enhancements.
- **Feedback**: Suggest UX tweaks or bold ideas.

---

## 📜 License

Licensed under the [MIT License](LICENSE). Hack freely, respect the code.

---

**"Stay low, move fast, control the network."**  
— NewCsTeamServer Crew