<template>
    <div class="spotify-player">
        <div id="player"></div>
        <button @click="play">Play</button>
    </div>
</template>

<script>
import axios from 'axios';

export default {
    name: "SpotifyPlayer",
    data() {
        return {
            player: null,
            deviceId: null,
            accessToken: null,
        };
    },
    async mounted() {
        try {
            const response = await axios.get("http://localhost:8080/token");
            this.accessToken = response.data.access_token;
            this.loadSpotifySDK(); // Now load the Spotify SDK after obtaining the token
        } catch (error) {
            console.error("Error getting access token:", error);
        }
    },
    methods: {
        loadSpotifySDK() {
            // Define the global function before loading the SDK
            window.onSpotifyWebPlaybackSDKReady = this.initializePlayer;

            // Dynamically load the Spotify SDK script
            const script = document.createElement('script');
            script.src = "https://sdk.scdn.co/spotify-player.js";
            script.async = true;
            script.onload = () => {
                console.log("Spotify Web Playback SDK script loaded successfully.");
            };
            script.onerror = () => {
                console.error("Failed to load the Spotify Web Playback SDK script.");
            };
            document.body.appendChild(script);
        },
        initializePlayer() {
            if (this.accessToken && window.Spotify) {
                console.log("initializing player");

                this.player = new Spotify.Player({
                    name: "Vue Spotify Player",
                    getOAuthToken: (cb) => {
                        cb(this.accessToken);
                    },
                    volume: 0.5,
                });

                this.player.on("initialization_error", (e) => console.error(e));
                this.player.on("account_error", (e) => console.error(e));
                this.player.on("playback_error", (e) => console.error(e));

                this.player.addListener("ready", ({ device_id }) => {
                    console.log("Ready with Device ID", device_id);
                    this.deviceId = device_id;
                });

                this.player.addListener("not_ready", ({ device_id }) => {
                    console.log("Device ID has gone offline", device_id);
                });

                this.player.connect().then((success) => {
                    if (success) {
                        console.log("Spotify player connected!");
                    }
                });
            } else {
                console.error("Spotify SDK not loaded properly or access token missing.");
            }
        },
        // play() {
        //     if (this.player) {
        //         this.player.resume().then(() => {
        //             console.log("Playback started");
        //         }).catch((error) => {
        //             console.error("Error starting playback:", error);
        //         });
        //     }
        // },
        play() {
            if (this.player) {
                axios.put(`https://api.spotify.com/v1/me/player/play?device_id=${this.deviceId}`, {
                    uris: ["spotify:track:6rqhFgbbKwnb9MLmUQDhG6"],
                }, {
                    headers: {
                        "Content-Type": "application/json",
                        Authorization: `Bearer ${this.accessToken}`,
                    },
                }).then(response => {
                    if (response.status === 200) {
                        console.log("Playback started");
                    } else {
                        throw new Error("Failed to start playback");
                    }
                }).catch(error => {
                    console.error("Error starting playback:", error);
                });

            }
        },
    },
};
</script>

<style scoped>
.spotify-player {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 100%;
}
</style>
