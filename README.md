# ThunderKitty
# DO NOT USE ITS NOT STABLE, WE WILL MODIFY THIS MESSAGE WHEN THUNDERKITTY IS IN STABLE STATE, CURRENTLY IT WILL BSOD YOU.
- 🔑 Open-source stealer written in Go, all logs will be sent to a Telegram bot.
- <a href="https://t.me/pulzetools"><img src="https://img.shields.io/badge/Join%20my%20Telegram%20group-2CA5E0?style=for-the-badge&logo=telegram&labelColor=db44ad&color=5e2775"></a>

![ThunderKitty Logo](https://github.com/EvilBytecode/ThunderKitty/assets/151552809/3c98bcf2-b958-48ae-8c5a-c8a0be13abd4)

---

## Current Features

- GUI Builder
- Anti-Kill: Terminating ThunderKitty will BSOD Computer
- Mutex (single instance)
- Antivirus Evasion:
  - Excludes current drive from Windows Defender
  - Lifetime AMSI and ETW patch
- Overtakes Hosts File
- Anti-Analysis for VMWare, VirtualBox, Sandboxes, Emulators, Debuggers, Any.run, Tria.ge
- Advanced AntiDebug
- Extracts WiFi Passwords
- Persistence:
  - Task Scheduler if Admin
  - Registry key if not Admin
- Backup Codes for Discord, Epic Games, Github
- Discord Token grabber for Discord PTB, Discord, Lightcord, Discord Canary
- Browsers:
  - Steals logins, cookies, credit cards, history, and download lists from 37 Chromium-based browsers
  - Steals logins, cookies, history, and download lists from 10 Gecko browsers
- Spreading: Executes and spreads specified message through Discord
- Disables TaskManager
- Disables Factory Reset
- Fake Error
- File Pumper
- Launches site
- Hides Console Window
- Swaps Mouse Buttons
- Changes Wallpaper
- Text-to-speech upon execution
- System Info

## Installation & Setup
### Pre-requisites:
- [The Go Programming Language](https://go.dev)
- [GCC/MinGW-w64](https://www.mingw-w64.org/)
- For help on installing MinGW-w64, consult [this link](https://code.visualstudio.com/docs/cpp/config-mingw).

### Building
- If you are building from source, you must run the following commands in your terminal.
- ```set CGO=1```
- ```go run main_gui.go```
- This command might take a few minutes as Go has to install packages such as Fyne, which are quite large.
- Once it finishes building, you will be presented with the builder UI and you will be able to proceed.

### Creating a Telegram Bot
- As this stealer uses Telegram for delivery of logs, you are required to create a bot.
- The first thing you must do is message [BotFather](https://t.me/botfather) and create a new bot.
- Once you have the bot created, message [chatIDrobot](https://t.me/chatIDrobot) to receive your chat ID.
- You can then put both of these values in the builder.
- Don't forget to start a conversation with the bot you just created, as you won't be able to receive messages otherwise.

---

## Detections

![Detection Image 1](https://github.com/EvilBytecode/ThunderKitty/assets/151552809/314a45d2-739f-4244-8daf-a257c61c133a)
![Detection Image 2](https://github.com/EvilBytecode/ThunderKitty/assets/151552809/0d773da7-3511-41e3-ac80-86dcf7b88f8d)
![Detection Image3](https://github.com/EvilBytecode/ThunderKitty/assets/151552809/93f2149c-099d-4af5-8f8d-e735db9c054e)

---

## GUI & Logs: 
![photo_2024-06-27_10-24-51](https://github.com/EvilBytecode/ThunderKitty/assets/151552809/493a1360-88eb-4cef-9ed6-11ea97c26354)
![photo_2024-06-27_10-24-51 (2)](https://github.com/EvilBytecode/ThunderKitty/assets/151552809/f85e40bc-cf49-465f-97e6-aedb8c829040)
![photo_2024-06-27_10-24-51 (3)](https://github.com/EvilBytecode/ThunderKitty/assets/151552809/ec1e5414-21dd-4cec-8585-17eeadc51060)
![image](https://github.com/EvilBytecode/ThunderKitty/assets/151552809/60c07839-33de-48ce-9db8-98f8f7a0bfbe)


---

## Credits

- [hackirby](https://github.com/hackirby) - Providing the base for the stealer.
- [SecDbg](https://github.com/secdbg) - Contributing heavily to development (Follow him, Deserved asf and hes 🐐).
- [KDot227](https://github.com/KDot227) - Hosts list.

 
![image](https://github.com/EvilBytecode/ThunderKitty/assets/151552809/09ce45b5-81d5-4940-a2d8-99706c5aaed1)


## FUTURE IDEAS FOR US
- Refracture Code (its messy)
- Games Stealing (Riot,Epic,Steam stealing)
- Self Destruct (Melt file)
- Kill Discord Token Protector
- Modify Assembly Info (Icon,Assm Info etc) (metadata)
- Documentation inside GUI, what each thing does.
- Encrypt Traffic
