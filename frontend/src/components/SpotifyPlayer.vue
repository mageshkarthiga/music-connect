<template>
    <div class="spotify-player">
        <Divider />
        <Paginator :rows="rowsPerPage" :totalRecords="totalTracks" :first="first" @page="onPageChange" />
        <Divider/>
        <div class="tracks">
            <Card v-for="track in paginatedTracks" :key="track.id"
                style="width: 25rem; overflow: hidden; margin-bottom: 20px;">
                <template #header>
                    <img :src="track.track_image_url" :alt="track.name" />
                </template>
                <template #title>{{ track.track_title }}</template>
                <template #footer>
                    <div class="flex gap-4 mt-1">
                        <Button icon="pi pi-caret-right" severity="success" rounded
                            @click="loadSpotifyContent(track.track_uri)" />
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
import axios from "axios";

export default {
    name: "SpotifyPlayer",
    data() {
        return {
            iframeSrc: "",
            accessToken: "",
            tracks: [], // All tracks fetched from the backend
            paginatedTracks: [], // Tracks to display on the current page
            totalTracks: 0, // Total number of tracks
            rowsPerPage: 12, // Number of tracks per page
            first: 0, // Index of the first track on the current page
        };
    },
    async mounted() {
        const response = await axios.get("http://localhost:8080/spotify/token");
        this.accessToken = response.data.access_token;
        this.loadSpotifyContent("spotify:track:5r7egnfTIQjaKSGREhIky9");

        await this.loadTracks();
        this.updatePaginatedTracks();
    },
    methods: {
        loadSpotifyContent(uri) {
            this.iframeSrc = `https://open.spotify.com/embed/${uri.split(":")[1]}/${uri.split(":")[2]}?autoplay=1`;
        },
        async loadTracks() {
            try {
                const response = await axios.get("http://localhost:8080/tracks");
                this.tracks = response.data;
                this.totalTracks = this.tracks.length; // Set total number of tracks
            } catch (error) {
                console.error("Error loading tracks:", error);
            }
        },
        updatePaginatedTracks() {
            // Calculate the tracks to display on the current page
            const start = this.first;
            const end = this.first + this.rowsPerPage;
            this.paginatedTracks = this.tracks.slice(start, end);
        },
        onPageChange(event) {
            // Update the `first` index when the page changes
            this.first = event.first;
            this.updatePaginatedTracks();
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
    grid-template-columns: repeat(auto-fit, minmax(20rem, 1fr));
    /* Reduced min width */
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