<template>
    <!-- Spotify Embed Iframe -->
    <div class="iframe-container">
        <iframe id="spotify-embed-iframe" :src="iframeSrc" frameborder="0" allowtransparency="true"
            allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture" loading="lazy"></iframe>
    </div>
</template>

<script>
import axios from "axios";
export default {
    name: "SpotifyPlayer",
    props: {
        spotifyUri: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            iframeSrc: "",
            accessToken: "",
        };
    },
    async mounted() {
        const response = await axios.get("http://localhost:8080/spotify/token");
        this.accessToken = response.data.access_token;
        this.loadSpotifyContent(this.spotifyUri);
    },
    watch: {
        spotifyUri(newUri) {
            this.loadSpotifyContent(newUri); 
        },
    },
    methods: {
        loadSpotifyContent(uri) {
            this.iframeSrc = `https://open.spotify.com/embed/${uri.split(":")[1]}/${uri.split(":")[2]}?autoplay=1`;
        },
    },
};
</script>

<style scoped>
.iframe-container {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    height: 80px;
    z-index: 1000;
    border-radius: 0;
}

.iframe-container iframe {
    width: 100%;
    height: 100%;
    border: none;
}
</style>