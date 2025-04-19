<template>
    <div class="music-tracks">
        <Divider />
        <Paginator :rows="rowsPerPage" :totalRecords="totalTracks" :first="first" @page="onPageChange" />
        <Divider />
        <div class="tracks">
            <Card
                v-for="track in paginatedTracks"
                :key="track.track_id"
                style="width: 20rem; overflow: hidden; margin-bottom: 20px;"
            >
                <template #header>
                    <img :src="track.track_image_url" :alt="track.track_title" />
                </template>
                <template #title>{{ track.track_title }}</template>
                <template #footer>
                    <div class="flex gap-4 mt-1">
                        <Button
                            icon="pi pi-play"
                            severity="success"
                            rounded
                            @click="() => { $emit('track-selected', track.track_uri); incrementPlayCount(track.track_id); }"
                        />
                    </div>
                </template>
            </Card>
        </div>
    </div>
</template>

<script>
import { incrementTrackPlayCount } from '@/service/TrackService';

export default {
    name: "MusicTracks",
    props: {
        tracks: {
            type: Array,
            required: true, // Ensure tracks are passed as a prop
        },
    },
    data() {
        return {
            localTracks: [], 
            paginatedTracks: [],
            totalTracks: 0,
            rowsPerPage: 12,
            first: 0,
        };
    },
    watch: {
        tracks: {
            immediate: true,
            handler(newTracks) {
                this.localTracks = [...newTracks]; 
                this.totalTracks = this.localTracks.length;
                this.updatePaginatedTracks(); 
            },
        },
    },
    methods: {
        updatePaginatedTracks() {
            const start = this.first;
            const end = this.first + this.rowsPerPage;
            this.paginatedTracks = this.localTracks.slice(start, end);
        },
        onPageChange(event) {
            this.first = event.first;
            this.updatePaginatedTracks();
        },
        incrementPlayCount(trackId) {
            incrementTrackPlayCount(trackId)
                .then(() => {
                    console.log(`Play count incremented for track ID: ${trackId}`);
                })
                .catch((error) => {
                    console.error(`Error incrementing play count for track ID: ${trackId}`, error);
                });
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
}

.track:hover {
    background: #1db954;
}
</style>