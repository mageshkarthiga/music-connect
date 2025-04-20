<template>
    <!-- Spotify Embed Iframe -->
    <div class="iframe-container">
        <iframe id="spotify-embed-iframe" :src="iframeSrc" frameborder="0" allowtransparency="true"
            allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture" loading="lazy"></iframe>
    </div>
</template>

<script>
import axios from "axios";
import { useSpotifyStore } from "@/store/SpotifyStore";

export default {
    name: "SpotifyPlayer",
    data() {
        return {
            iframeSrc: "",
            accessToken: "",
            API_BASE_URL: process.env.VUE_APP_API_BASE_URL,
        };
    },
    async mounted() {
        const spotifyStore = useSpotifyStore();
        const response = await axios.get(`${this.API_BASE_URL}/spotify/token`);
        this.accessToken = response.data.access_token;
        this.loadSpotifyContent(spotifyStore.spotifyUri);

        // Watch for store updates
        this.$watch(
            () => spotifyStore.spotifyUri,
            (newUri) => {
                this.loadSpotifyContent(newUri);
            }
        );
    },
    methods: {
        loadSpotifyContent(uri) {
            if (uri && uri.includes(":")) {
                const [type, id] = uri.split(":").slice(1);
                this.iframeSrc = `https://open.spotify.com/embed/${type}/${id}?autoplay=1`;
            }
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
