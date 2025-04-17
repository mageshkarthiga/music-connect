from flask import Flask, request, jsonify
from scraper import scrape_ticketmaster
import threading

app = Flask(__name__)

@app.post("/scrape/ticketmaster")
def ticketmaster_scrape():
    data = request.get_json()
    callback_url = data.get("callbackUrl")

    if not callback_url:
        return jsonify({"error": "Missing callback_url"}), 400

    threading.Thread(target=scrape_ticketmaster, args=(callback_url,)).start()
    return jsonify({"status": "Scraping started"}), 200

if __name__ == "__main__":
    app.run(debug=True, port=3001)
