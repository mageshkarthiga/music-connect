<template>
    <AppTopbar />
    <div className="card mt-10">
        <h2 class="music-title">Spotify Music Discovery 🎶🎧</h2>
        <SpotifyPlayer />
    </div>
</template>

<script>
import SpotifyPlayer from '@/components/SpotifyPlayer.vue';
import AppTopbar from "@/layout/AppTopbar.vue"; // Adjust the path if necessary


export default {
    components: {
        SpotifyPlayer,
        AppTopbar,
    },
};
</script>

<style scoped>
.music-title {
    font-size: 2rem;
    font-weight: bold;
    margin-bottom: 1rem;
    text-align: center;
    position:sticky;
    top:0;
}

</style>