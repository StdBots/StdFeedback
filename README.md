# STD StdFeedback - Advanced Telegram Feedback Bot

A modern, highly-customizable Telegram Feedback & Support bot written in Go (Golang).

## Features
- Complete Ticket Management System (Open, Close, Resolve)
- Forwarding and Copy Broadcasting
- Custom Start Texts and Donation Links
- Role-based Access (Owner vs Auth Users)
- Inline Keyboards and Dynamic Responses
- Ban/Unban Users with reasons and durations
- Real-time MongoDB integrations
- Detailed Statistics Tracking (Response times, daily active users)
- Sub-bot Cloning feature (Create multiple bots per user)
- Logging to dedicated channels
- Webhooks and Polling support
- Anti-spam & Rate limiting
- Labels/Tags for support tickets
- User ratings upon ticket resolution
- Completely Open Source (AGPL-3.0)

## Quick Deploy Guide
1. **Clone the repository:**
   ```bash
   git clone https://github.com/StdBots/StdFeedback.git
   cd StdFeedback
   ```
2. **Configure Environment:**
   Copy `.env.example` to `.env` and fill the variables.
3. **Run using Docker:**
   ```bash
   docker build -t feedvackbot .
   docker run --env-file .env -d feedvackbot
   ```
   Or using Go directly:
   ```bash
   go run cmd/bot/main.go
   ```

## License & Credits
This project is licensed under the **AGPL-3.0 License**.

**CREDIT REQUIREMENTS:**
As per the AGPL-3.0 license and project terms, any forks, clones, or deployed instances of this bot MUST retain credits to the original developer.
- Developer: **STD DEEPANSHU** ([Website](https://deepanshu.in))
- Telegram Channels: [@STDBOTS](https://t.me/STDBOTS) | [@STD_DEEPANSHU](https://t.me/STD_DEEPANSHU)

If you modify or run this bot, the start menu or help section must contain clear links referencing STD DEEPANSHU and STD BOTS.

---
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)
![Telegram](https://img.shields.io/badge/Telegram-2CA5E0?style=for-the-badge&logo=telegram&logoColor=white)
