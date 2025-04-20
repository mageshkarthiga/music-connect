from flask import Flask, request, jsonify
from flask_cors import CORS
from scraper import scrape_ticketmaster
import threading
import os

app = Flask(__name__)
CORS(app)

@app.post("/scrape/ticketmaster")
def ticketmaster_scrape():
    data = request.get_json()
    callback_url = data.get("callbackUrl")

    if not callback_url:
        return jsonify({"error": "Missing callback_url"}), 400

    threading.Thread(target=scrape_ticketmaster, args=(callback_url,)).start()
    return jsonify({"status": "Scraping started"}), 200

if __name__ == "__main__":
    port = int(os.environ.get("PORT", 3001))
    app.run(host='0.0.0.0', debug=True, port=port)
