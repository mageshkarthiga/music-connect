from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from webdriver_manager.chrome import ChromeDriverManager
from bs4 import BeautifulSoup
import requests
import time

BASE_URL = "https://ticketmaster.sg"

def get_driver():
    chrome_options = Options()
    chrome_options.add_argument("--headless")
    chrome_options.add_argument("--no-sandbox")
    chrome_options.add_argument("--disable-dev-shm-usage")
    chrome_options.add_argument("user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/124.0.0.0 Safari/537.36")

    service = Service(ChromeDriverManager().install())
    return webdriver.Chrome(service=service, options=chrome_options)

def get_event_links(driver):
    driver.get(BASE_URL)
    time.sleep(3)
    soup = BeautifulSoup(driver.page_source, "html.parser")

    event_links = []
    for a in soup.select("a[href*='/activity/detail/']"):
        href = a.get("href")
        if href and href not in event_links:
            event_links.append(BASE_URL + href)
    return event_links

def extract_event_data(driver, url):
    driver.get(url)
    time.sleep(2)
    soup = BeautifulSoup(driver.page_source, "html.parser")

    try:
        event_name = soup.find("h1").get_text(strip=True)
    except:
        event_name = "Unknown Event"

    try:
        desc_tag = soup.find("meta", {"name": "description"})
        event_description = desc_tag["content"] if desc_tag else "No description"
    except:
        event_description = "No description"

    try:
        img_tag = soup.find("meta", {"name": "og:image"})
        image_url = image_url = img_tag["content"] if img_tag else ""
    except:
        image_url = ""

    try:
        venue_tag = soup.find("span", {"id": "synopsisEventVenue"})
        venue_name = venue_tag.get_text(strip=True) if venue_tag else "Unknown Venue"
    except:
        venue_name = "Unknown Venue"

    return {
        "eventName": event_name,
        "eventDescription": event_description,
        "eventUrl": url,
        "eventImageUrl": image_url,
        "venueName": venue_name,
        "location": "Singapore"
    }

def scrape_ticketmaster(callback_url):
    driver = get_driver()
    try:
        event_urls = get_event_links(driver)
        print(f"Found {len(event_urls)} events.")

        all_data = []
        for url in event_urls:
            print(f"Scraping: {url}")
            try:
                data = extract_event_data(driver, url)
                all_data.append(data)
            except Exception as e:
                print(f"Failed to scrape {url}: {e}")
            time.sleep(1)

        print(f"Total events scraped: {len(all_data)}")
        print("Example: ", all_data[0])
    finally:
        driver.quit()
    
    try:
        res = requests.post(callback_url, json=all_data)
        print(f"Callback status: {res.status_code}")
    except Exception as e:
        print(f"Error posting to callback: {e}")
