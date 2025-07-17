# 🚀 Torrent-Hoarding

**Automate torrent downloads and upload them directly to your favorite cloud storages!**

---

## 🌐 Supported Cloud Uploads

| Platform | Status | Features |
|----------|--------|----------|
| 📱 **Telegram** | ✅ **Fully Supported** | Auto thumbnails, Large file support |
| 🗂️ **Google Drive** | 🔜 **Coming Soon** |  |
| 🎬 **Streamtape** | 🔜 **Coming Soon** | |

---

## 📦 Telegram Features

### 📸 **Auto Thumbnails & Screenshots**
- Automatically generate video thumbnails and screenshots using FFmpeg
- Professional-looking file previews
- Configurable thumbnail quality

### 📤 **Large File Upload Support**
- 📁 **Free Users**: Upload up to **2GB** per file
- 💎 **Telegram Premium**: Upload up to **4GB** per file
- Smart file splitting for larger content

---

## ✨ Telegram Setup Guide

Setting up Telegram with Torrent-Hoarding is quick and easy! Just follow these steps:

### 1. 🔑 Create a Telegram App
1. Visit [my.telegram.org](https://my.telegram.org)
2. Log in with your phone number
3. Go to **API Development Tools**
4. Create a new application and copy your:
   - `API_ID`
   - `API_HASH`

### 2. 📢 Create a Telegram Channel or Group
1. Create a private **Channel** or **Group** where uploads will go
2. Use [@userinfobot](https://t.me/userinfobot) to get the **Chat ID** of your channel or group

### 3. ⚙️ Configure the App
1. Open the config file
2. Fill in the required fields:
   - `API_ID`, `API_HASH`
   - Chat/Channel ID
   - Session Name
   - Download Path

> 💡 **Pro Tip**: Make sure your credentials are correct to avoid upload issues!

---

## ⚡ Key Features

### 🧲 **Magnet Link Support Only**
- Add torrents using magnet links for simplicity and speed
- No need to download `.torrent` files

### 📥 **Effortless Torrent Downloads**
- CLI-based downloading — no bloated UI!

### 🎞️ **Automatic Video Thumbnails**
- Auto-generate thumbnails using ffmpeg (optional)
- Perfect for video content organization

### 🧹 **Smart Cleanup (Configurable)**
- Automatically delete local files after upload
- Free up disk space efficiently

---

## 🛠️ Installation & Setup

### 1. ✅ Prerequisites
- **Go** (version 1.18 or above recommended)
- Clone this repository:

```bash
git clone https://github.com/yourusername/torrent-hoarding.git
cd torrent-hoarding
```

### 2. 🧹 Install Dependencies
```bash
go mod tidy
```

### 3. ⚙️ Build the Application
```bash
go build .
```

### 4. 🚀 Run the App
```bash
./torrent-hoarding
```

---

## 🎬 FFmpeg Setup (Optional)

FFmpeg is required for video thumbnails and screenshots. Skip this if you don't need that feature.

### 🪟 Windows
1. Download FFmpeg from [ffmpeg.org](https://ffmpeg.org/download.html)
2. Extract the ZIP file
3. Add the `bin` folder to your System Environment Variables → Path
4. Test installation:
```bash
ffmpeg -version
```

### 🐧 Linux (Debian/Ubuntu)
```bash
sudo apt update
sudo apt install ffmpeg
```

### 🧊 NixOS
**For temporary use:**
```bash
nix-shell -p ffmpeg
```

**For permanent setup**, add to your `configuration.nix`:
```nix
environment.systemPackages = with pkgs; [
  ffmpeg
];
```

Then run:
```bash
sudo nixos-rebuild switch
```


---

## 📋 System Requirements

- **OS**: Linux, Windows, macOS
- **Go**: Version 1.18+

---

## 🤝 Contributing

We welcome contributions! Please feel free to submit a Pull Request.

---

## ⚠️ Important Notes

> **Cloud Engine Optimization**: This tool is designed for cloud engine usage. Concurrent downloads are not yet supported to prevent storage overuse and avoid upload bans.

> **Disclaimer**: This project is not affiliated with Telegram or any other cloud storage provider. You are solely responsible for how you use this tool. Use at your own risk.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🌟 Support

If you find this project helpful, please consider giving it a ⭐ on GitHub!

**Happy Torrenting!** 🎉