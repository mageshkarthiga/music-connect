<template>
    <AppTopbar />
    <div class="card mt-10">
        <h2 class="music-title">{{ playlist_name }} </h2>
        <!-- Pass tracks as a prop to MusicTracks -->
        <MusicTracks :tracks="tracks" @track-selected="onTrackSelected" />
        <!-- SpotifyPlayer Component -->
        <SpotifyPlayer />
    </div>
</template>

<script>
import MusicTracks from '@/components/MusicTracks.vue';
import SpotifyPlayer from '@/components/SpotifyPlayer.vue';
import { useSpotifyStore } from '@/store/SpotifyStore';
import axios from 'axios';

export default {
    components: {
        SpotifyPlayer,
        MusicTracks,
    },
    props: {
        playlist_id: {
            type: String,
            required: true,
        },
        playlist_name: {
            type: String,
            required: true,
        },
    },
    data() {
        return {
            tracks: [], 
        };
    },
    async mounted() {
        await this.fetchPlaylistTracks();
    },
    watch: {
        
        playlist_id: "fetchPlaylistTracks",
        playlist_name: "fetchPlaylistTracks",
    },
    methods: {
        async fetchPlaylistTracks() {
            try {
                const response = await axios.get(`http://localhost:8080/playlists/${this.playlist_id}/tracks`);
                this.tracks = response.data; // Populate tracks
            } catch (error) {
                console.error("Error fetching playlist tracks:", error);
            }
        },
        onTrackSelected(uri) {
            const spotifyStore = useSpotifyStore();
            spotifyStore.setSpotifyUri(uri);
        },
    },
};
</script>

<style scoped>
.music-title {
    font-size: 2rem;
    font-weight: bold;
    margin-bottom: 1rem;
    text-align: center;
    position: sticky;
    top: 0;
}
</style>