# Use official Python image
FROM python:3.10-slim

# Prevents interactive prompts during build
ENV DEBIAN_FRONTEND=noninteractive

# Install Chrome browser and necessary tools
RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    wget \
    gnupg \
    libnss3 \
    libgconf-2-4 \
    libxss1 \
    libappindicator1 \
    libindicator7 \
    fonts-liberation \
    xdg-utils \
    libu2f-udev \
    ca-certificates \
    libasound2 \
    libatk-bridge2.0-0 \
    libatk1.0-0 \
    libcups2 \
    libdbus-1-3 \
    libgtk-3-0 \
    libx11-xcb1 \
    libxcomposite1 \
    libxdamage1 \
    libxrandr2 \
    x11-utils \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Download and install Chrome (Stable)
RUN wget -q https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    dpkg -i google-chrome-stable_current_amd64.deb || apt-get -fy install && \
    rm google-chrome-stable_current_amd64.deb

# Set working directory
WORKDIR /app

# Copy requirements and install Python dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy your app files
COPY . .

# Expose port for Flask
EXPOSE 3001

# Run Flask app
CMD ["python", "app.py"]
