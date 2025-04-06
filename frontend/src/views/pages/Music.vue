<template>
    <div className="card">
        <h2 className="text-2xl font-bold mb-4">Spotify Music Discovery</h2>
        <SpotifyPlayer />
    </div>
</template>

<script>
import SpotifyPlayer from '@/components/SpotifyPlayer.vue';


export default {
    components: {
        SpotifyPlayer,
    },
};
</script>