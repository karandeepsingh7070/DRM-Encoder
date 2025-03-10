🎬 DRM Video Encoder & Player

📌 Overview

DRM Video Encoder & Player is a web application that allows users to convert normal videos to DRM-encrypted formats (Widevine, PlayReady, Raw Key) and test playback using Shaka Player. The backend is built using Go and Shaka Packager, while the frontend is developed with Next.js and Tailwind CSS.

✨ Features

Upload MP4 videos and select encryption type (Widevine, PlayReady, or Raw Key).

Choose segment size and audio inclusion.

Download the encrypted files after processing.

Shaka Player integration for testing DRM playback.

🏗️ Tech Stack

Frontend (Next.js + Tailwind CSS)

Next.js 
Tailwind CSS 
Shaka Player
Backend (Go + Shaka Packager)
FFmpeg (if required for pre-processing)

🚀 Installation & Setup

1️⃣ Clone the Repository

git clone https://github.com/your-username/drm-encoder.git
cd drm-encoder

2️⃣ Backend Setup (Go Server)

Install Dependencies

Ensure you have Go and Shaka Packager installed:

# Install Go
https://go.dev/doc/install

# Install Shaka Packager
https://github.com/shaka-project/shaka-packager

Run Backend Server

cd backend
go mod tidy
go run main.go

The backend runs on http://localhost:8080.

3️⃣ Frontend Setup (Next.js)

Install Dependencies

cd client
yarn install  # or npm install

Run Frontend

yarn dev  # or npm run dev

The frontend runs on http://localhost:3000.

🎥 Usage Guide

1️⃣ Upload & Encrypt Video

Select a video file.

Choose encryption type (Widevine, PlayReady, or Raw Key).

Adjust segment size and enable/disable audio.

Click Convert Video to start encryption.

Download encrypted files after processing.

2️⃣ Test DRM Playback

Switch to the Test Playback tab.

Provide the MPD URL from the encrypted files.

Play video using Shaka Player.

🤝 Contributing

Pull requests and feature suggestions are welcome! 🚀

For any queries, feel free to reach out!