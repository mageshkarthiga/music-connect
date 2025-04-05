<template>
    <div class="spotify-player">
        <div class="tracks">
            <Card v-for="track in tracks" style="width: 25rem; overflow: hidden; margin-bottom: 20px;"
                :key="track.id" :data-spotify-id="track.uri">
                <template #header>
                    <img :src="track.image" :alt="track.name" />
                </template>
                <template #title>{{ track.name }}</template>
                <template #subtitle>{{ track.album }}</template>
                <template #content>
                    <p class="track-description">
                        {{ truncateDescription(track.description) }}
                    </p>
                </template>
                <template #footer>
                    <div class="flex gap-4 mt-1">
                        <Button icon="pi pi-caret-right" severity="success" rounded @click="loadSpotifyContent(track.uri)"/>
                    </div>
                </template>
            </Card>
        </div>

        <!-- Spotify Embed Iframe -->
        <div class="iframe-container">
            <iframe id="spotify-embed-iframe" :src="iframeSrc" frameborder="0" allowtransparency="true"
                allow="autoplay; clipboard-write; encrypted-media; fullscreen; picture-in-picture"
                loading="lazy"></iframe>
        </div>
    </div>

</template>

<script>
// import { useSpotifyStore } from "@/stores/spotifyStore";
import axios from "axios";

export default {
    name: "SpotifyPlayer",
    data() {
        return {
            iframeSrc: "",
            accessToken: "",
            tracks: [
            ],
        };
    },
    async mounted() {
        const response = await axios.get("http://localhost:8080/spotify/token");
        this.accessToken = response.data.access_token;
        this.loadSpotifyContent("spotify:track:5r7egnfTIQjaKSGREhIky9");

        this.tracks = this.loadtracks();
        if (tracks) {
            this.tracks = tracks;
        } else {
            console.error("Failed to load tracks");
        }
    },
    methods: {
        loadSpotifyContent(uri) {
            // const spotifyStore = useSpotifyStore();

            this.iframeSrc = `https://open.spotify.com/embed/${uri.split(":")[1]}/${uri.split(":")[2]}?autoplay=1`;

            // spotifyStore.setCurrentUri(uri);
        },
        loadtracks() {
            const apiUrl = "https://api.spotify.com/v1/tracks"
            axios.get(apiUrl, {
                headers: {
                    Authorization: `Bearer ${this.accessToken}`,
                },
                params: {
                    ids: "5biPEkF6O8snl96CWzsdK3,6GygzQ98i5iORNo45s5KwF,5r7egnfTIQjaKSGREhIky9,0fK7ie6XwGxQTIkpFoWkd1,0b0Dz0Gi86SVdBxYeiQcCP",
                },
            })
                .then((response) => {
                    this.tracks = response.data.tracks.map((track) => ({
                        id: track.id,
                        name: track.name,
                        album: track.album.name,
                        uri: track.uri,
                        image: track.album.images[0]?.url || "", 
                        description: `${track.artists.map((artist) => artist.name).join(", ")}`,
                    }));
                    console.log("tracks loaded:", this.tracks);
                })
                .catch((error) => {
                    console.error("Error fetching tracks:", error);
                });
        },
        truncateDescription(description) {
            const maxLength = 300; // Set the maximum length for the description
            if (description.length > maxLength) {
                return description.substring(0, maxLength) + "...";
            }
            return description;
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

.tracks {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(20rem, 1fr)); /* Reduced min width */
    gap: 15px;
    margin-bottom: 1rem;
    width: 100%;
}

.track {
    min-width: max-content;
    margin-bottom: 0.8rem;
    padding: 0.8rem 1rem;
    border-radius: 10px;
    border: 0;
    background: #191414;
    color: #fff;
    cursor: pointer;
}

.track:hover {
    background: #1db954;
}

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