# 🎵 Music Connect

**Music Connect** is a social music platform that allows users to discover what people nearby are listening to, connect through real-time messaging, match with music lovers, and receive personalised event and music recommendations. The platform combines music discovery, social interaction, and event management into a seamless, location-based experience.

---

## 📑 Table of Contents

- [✨ Features](#-features)
- [🚀 Frontend Setup](#-frontend-setup)
- [🔧 Backend Setup](#-backend-setup)
  - [📍 Location Autofill Service](#-location-autofill-service)
- [▶️ Usage](#-usage)
- [🪪 License](#-license)

---

## ✨ Features

- 🤝 **Friend Matching & Messaging**  
  Connect with users who share your music taste and chat with them through the in-app messaging system.

- 🔍 **Music Discovery & Recommendations**  
  Discover what nearby users are listening to and get personalised music suggestions based on your and your friends’ playlists.

- 📍 **Nearby User Locator**  
  Find and interact with users around your current location.

- 🎶 **Listen, Like & Playlist Management**  
  Play music, like tracks, and curate your own playlists with ease.

- 🎫 **Event Suggestions**  
  Receive personalised event recommendations based on your music taste and social activity.

---

## 🚀 Frontend Setup

To run the frontend locally:

1. Navigate to the project directory:

   ```bash
   cd frontend/src
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Start the development server:

   ```bash
   npm run dev
   ```

4. Open your browser and go to:
   ```
   http://localhost:5173/
   ```

---

## 🔧 Backend Setup

To run the backend services:

1. Navigate to the backend directory:

   ```bash
   cd backend
   ```

2. Set up each backend service as per their individual setup guides (e.g. authentication, chat service, music recommender, etc.).

---

### 📍 Location Autofill Service

This service is written in Go and provides location autocomplete functionality:

1. Go to the `map` service directory:

   ```bash
   cd map
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Run the service:
   ```bash
   go run main.go
   ```

---

## ▶️ Usage

1. Start the frontend server ([see steps](#frontend-setup)).
2. Start all required backend services ([see steps](#backend-setup)).
3. Open the app in your browser and dive into the music-powered social world!

---

## 🪪 License

This project is licensed under the [MIT License](LICENSE). See the `LICENSE` file for more details.
