# Music Connect

**Music Connect** is a social music platform that allows users to discover what people nearby are listening to, connect through real-time messaging, match with music lovers, and receive personalised event and music recommendations. The platform combines music discovery, social interaction, and event management into a seamless, location-based experience.

---

## Table of Contents

- [Features](#features)
- [Frontend Setup](#frontend-setup)
- [Backend Setup](#backend-setup)
- [Usage](#usage)

---

## Features

- **Friend Matching & Messaging**  
  Connect with users who share your music taste and chat with them through the in-app messaging system.

- **Music Discovery & Recommendations**  
  Discover what nearby users are listening to and get personalised music suggestions based on your and your friends’ playlists.

- **Nearby User Locator**  
  Find and interact with users around your current location.

- **Listen, Like & Playlist Management**  
  Play music, like tracks, and curate your own playlists with ease.

- **Event Suggestions**  
  Receive personalised event recommendations based on your music taste and social activity.

---

## Frontend Setup

To run the frontend locally:

1. Navigate to the project directory:

   ```bash
   cd frontend
   ```

2. Install dependencies:

   ```bash
   npm install
   ```

3. Start the development server:

   ```bash
   npm run dev
   ```

4. Open your browser and go to: [http://localhost:5173/](http://localhost:5173/)

---

## Backend Setup

To run the backend services:

1. Navigate to the backend directory:

   ```bash
   cd backend
   ```

2. Install dependencies:

   ```go
   go mod tidy
   ```

3. Build and run the service:

   ```go
   go build
   go run .
   ```

---
### Chat Service

This service is written in Go and provides real-time messaging functionality. It allows users to:
- Send and receive messages from friends in real time.
- Store chat history.

To set up the chat service:

1. Navigate to the `chat-service` directory:

   ```bash
   cd chat-service
   ```

2. Install dependencies:

   ```go
   go mod tidy
   ```

3. Run the service:

   ```go
   go run main.go
   ```

---

### Location Service Setup

This service is written in Go and provides location autocomplete functionality:

1. Go to the `map` service directory:

   ```bash
   cd map
   ```

2. Install dependencies:

   ```go
   go mod tidy
   ```

3. Run the service:

   ```go
   go run main.go
   ```

---

### Cron Job to populate events and tracks

Two services are used in this system: one written in Go and the other in Python.
The Go service includes a scheduled job that runs every Sunday at 00:00. This job performs two tasks:

- Calls the Spotify API to check for and store any newly released tracks.

- Triggers the Python Scraper service, which scrapes Ticketmaster for new events and returns the data to be stored.

To setup:

1. Open 2 terminals. In the first terminal, go to the `ticketmaster scraper` service directory:

   ```bash
   cd ticketmaster-scraper
   ```

2. Run the service:

   ```bash
   python3 app.py
   ```

3. In the second terminal, go to the `cron-jobs` service directory:

   ```bash
   cd cron-jobs
   ```

4. Install dependencies:

   ```go
   go mod tidy
   ```

5. Build and run the service:

   ```go
   go build
   go run .
   ```

---

## Usage

1. Start the frontend server ([see steps](#frontend-setup)).
2. Start all required backend services ([see steps](#backend-setup)).
3. Open the app in your browser and dive into the music-powered social world!
