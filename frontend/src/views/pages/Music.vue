<template>
    <AppTopbar />
    <div className="card mt-10">
        <h2 class="music-title">Your Tracks 🎶🎧</h2>
        <MusicTracks @track-selected="onTrackSelected" />
        <!-- SpotifyPlayer Component -->
        <SpotifyPlayer v-if="selectedTrackUri" :spotifyUri="selectedTrackUri" />
    </div>
</template>

<script>
import MusicTracks from '@/components/MusicTracks.vue';
import SpotifyPlayer from '@/components/SpotifyPlayer.vue';


export default {
    components: {
        SpotifyPlayer,
        MusicTracks,
    },
    data() {
        return {
            selectedTrackUri: "spotify:track:5r7egnfTIQjaKSGREhIky9", 
        };
    },
    methods: {
        onTrackSelected(uri) {
            this.selectedTrackUri = uri; 
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