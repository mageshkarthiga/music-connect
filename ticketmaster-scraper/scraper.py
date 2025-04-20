from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from bs4 import BeautifulSoup
from webdriver_manager.chrome import ChromeDriverManager
import requests
import time
import os

BASE_URL = "https://ticketmaster.sg/"

def get_driver():
    print("Initializing WebDriver...")
    chrome_options = Options()
    chrome_options.add_argument("--headless")
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument("--disable-dev-shm-usage")

    # Check if running in Docker or locally
    environment = os.getenv("ENVIRONMENT", "local")
    print(f"Running in {environment} environment")

    if environment == "docker":
        chrome_options.add_argument("--disable-gpu")
        chrome_options.add_argument("--remote-debugging-port=9222")
        chrome_options.binary_location = "/usr/bin/chromium"
        # Use paths for Docker
        chrome_options.binary_location = "/usr/bin/chromium"
        service = Service(executable_path="/usr/bin/chromedriver")
    else:
        service = Service(ChromeDriverManager().install())

    driver = webdriver.Chrome(service=service, options=chrome_options)
    print("WebDriver initialized successfully.")

    return driver

def get_event_links(driver):
    print("Opening the base URL...")
    driver.get(BASE_URL)
    time.sleep(3)
    soup = BeautifulSoup(driver.page_source, "html.parser")
    print("Page loaded. Extracting event links...")

    event_links = []
    for a in soup.select("a[href*='/activity/detail/']"):
        href = a.get("href")
        if href and href not in event_links:
            event_links.append(BASE_URL + href)
    print(f"Found {len(event_links)} event links.")
    return event_links

def extract_event_data(driver, url):
    print(f"Extracting data from event: {url}")
    driver.get(url)
    time.sleep(2)
    soup = BeautifulSoup(driver.page_source, "html.parser")

    try:
        event_name = soup.find("h1").get_text(strip=True)
        print(f"Event name found: {event_name}")
    except:
        event_name = "Unknown Event"
        print("Event name not found.")

    try:
        desc_tag = soup.find("meta", {"name": "description"})
        event_description = desc_tag["content"] if desc_tag else "No description"
        print(f"Event description: {event_description}")
    except:
        event_description = "No description"
        print("Event description not found.")

    try:
        img_tag = soup.find("meta", {"name": "og:image"})
        image_url = img_tag["content"] if img_tag else ""
        print(f"Image URL: {image_url}")
    except:
        image_url = ""
        print("Image URL not found.")

    try:
        venue_tag = soup.find("span", {"id": "synopsisEventVenue"})
        venue_name = venue_tag.get_text(strip=True) if venue_tag else "Unknown Venue"
        print(f"Venue: {venue_name}")
    except:
        venue_name = "Unknown Venue"
        print("Venue not found.")

    return {
        "eventName": event_name,
        "eventDescription": event_description,
        "eventUrl": url,
        "eventImageUrl": image_url,
        "venueName": venue_name,
        "location": "Singapore"
    }

def scrape_ticketmaster(callback_url):
    print("Starting scraping process...")
    driver = get_driver()
    try:
        print("Getting event links...")
        event_urls = get_event_links(driver)
        print(f"Found {len(event_urls)} events.")

        all_data = []
        for url in event_urls:
            print(f"Scraping event: {url}")
            try:
                data = extract_event_data(driver, url)
                all_data.append(data)
                print(f"Scraped data for event: {url}")
            except Exception as e:
                print(f"Failed to scrape {url}: {e}")
            time.sleep(1)

        print(f"Total events scraped: {len(all_data)}")
        print(f"Example event data: {all_data[0] if all_data else 'No data'}")
    finally:
        driver.quit()
        print("WebDriver closed.")
    
    try:
        print("Sending scraped data to callback URL...")
        res = requests.post(callback_url, json=all_data)
        print(f"Callback status: {res.status_code}")
    except Exception as e:
        print(f"Error posting to callback: {e}")

