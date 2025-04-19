# Music Connect

Music Connect is a web application that allows users to discover what people nearby are listening to, connect through messaging, and receive personalized event recommendations. The platform combines music discovery, social interaction, and event management into a seamless experience.

---

## Table of Contents

- [Features](#features)
- [Frontend Setup](#frontend-setup)
- [Backend Setup](#backend-setup)
  - [Location Autofill Service](#location-autofill-service)
- [Usage](#usage)
- [Contributing](#contributing)
- [License](#license)

---

## Features

- **Music Discovery**: See what music people nearby are listening to.
- **Messaging**: Connect with others through an integrated chat system.
- **Event Recommendations**: Get personalized event suggestions based on your music preferences.
- **Responsive Design**: Optimized for both desktop and mobile devices.

---

## Frontend Setup

To set up the frontend locally:

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

4. Open your browser and navigate to:
   ```
   http://localhost:5173/
   ```

---

## Backend Setup

To set up the backend:

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Follow the specific instructions for each backend service.

---

### Location Autofill Service

The location autofill service is written in Go. To set it up:

1. Navigate to the `map` directory:
   ```bash
   cd map
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the service:
   ```bash
   go run filename.go
   ```

---

## Usage

1. Start the frontend server using the steps in [Frontend Setup](#frontend-setup).
2. Start the backend services using the steps in [Backend Setup](#backend-setup).
3. Open the application in your browser and explore its features.

---

## License

This project is licensed under the [MIT License](LICENSE). See the `LICENSE` file for details.
