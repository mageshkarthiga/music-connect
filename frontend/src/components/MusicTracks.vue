<template>
    <div class="music-tracks">
        <Divider />
        <Paginator :rows="rowsPerPage" :totalRecords="totalTracks" :first="first" @page="onPageChange" />
        <Divider/>
        <div class="tracks">
            <Card v-for="track in paginatedTracks" :key="track.id"
                style="width: 20rem; overflow: hidden; margin-bottom: 20px;">
                <template #header>
                    <img :src="track.track_image_url" :alt="track.name" />
                </template>
                <template #title>{{ track.track_title }}</template>
                <template #footer>
                    <div class="flex gap-4 mt-1">
                        <Button icon="pi pi-play" severity="success" rounded
                            @click="$emit('track-selected', track.track_uri)" />
                    </div>
                </template>
            </Card>
        </div>
    </div>
</template>

<script>
import axios from "axios";
import SpotifyPlayer from "./SpotifyPlayer.vue";

export default {
    name: "MusicTracks",
    components: {
        SpotifyPlayer,
    },
    data() {
        return {
            iframeSrc: "",
            accessToken: "",
            tracks: [], 
            paginatedTracks: [],
            totalTracks: 0,
            rowsPerPage: 12,
            first: 0,
        };
    },
    async mounted() {
        const response = await axios.get("http://localhost:8080/spotify/token");
        this.accessToken = response.data.access_token;

        await this.loadTracks();
        this.updatePaginatedTracks();
    },
    methods: {
        selectTrack(uri) {
            this.selectedTrackUri = uri; 
        },
        async loadTracks() {
            try {
                const response = await axios.get("http://localhost:8080/me/tracks",{
                    withCredentials: true,
                });
                this.tracks = response.data;
                this.totalTracks = this.tracks.length; 
            } catch (error) {
                console.error("Error loading tracks:", error);
            }
        },
        updatePaginatedTracks() {
            const start = this.first;
            const end = this.first + this.rowsPerPage;
            this.paginatedTracks = this.tracks.slice(start, end);
        },
        onPageChange(event) {
            this.first = event.first;
            this.updatePaginatedTracks();
        },
    },
};
</script>

<style scoped>
.music-tracks {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    height: 100%;
}

.tracks {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(20rem, 1fr));
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
</style>