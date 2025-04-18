# Use official Python slim image
FROM python:3.10-slim

# Avoid interactive prompts during package installs
ENV DEBIAN_FRONTEND=noninteractive

# Install Chromium and Chromium driver for headless browsing
RUN apt-get update && apt-get install -y \
    chromium \
    chromium-driver \
    fonts-liberation \
    wget \
    curl \
    gnupg \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy app code and requirements
COPY . .

# Install Python dependencies
RUN pip install --no-cache-dir -r requirements.txt

# Expose port (for Flask or FastAPI etc.)
EXPOSE 8080

# Default command to run your app
CMD ["python", "app.py"]
