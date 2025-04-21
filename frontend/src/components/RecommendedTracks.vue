<script>
import { getRecommendedTracks } from '@/service/RecommenderService';
import trackService from '@/service/TrackService'; // Import track service to handle API calls for liking and unliking tracks

export default {
    name: 'RecommendedTracks',
    data() {
        return {
            tracks: null,
            selectedTrackUri: null,
            likedTracks: new Set(), // to store liked track IDs
        };
    },
    methods: {
        async fetchRecommendedTracks() {
            try {
                const data = await getRecommendedTracks();
                this.tracks = data;

                console.log('Recommended tracks:', this.tracks);
            } catch (error) {
                console.error('Error fetching recommended tracks:', error);
            }
        },
        handleClick(trackUri) {
            this.selectedTrackUri = trackUri;
            this.$emit('track-selected', trackUri);
        },

        async likedTracksHandler() {
            try {
        
                const liked = await trackService.likedTracks(); 

                if (!Array.isArray(liked)) {
                    throw new Error('Invalid liked tracks data format');
                }

                this.likedTracks = new Set(liked.map(t => t.track_id));
            } catch (error) {
                console.error('Error initializing tracks:', error);
            } finally {
                this.loading = false;
            }
        },

        async likeTrackHandler(trackId) {
            try {
                // Like the track locally and on the server
                this.likedTracks.add(trackId);
                this.$toast.add({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'Track liked successfully!',
                });
                await trackService.likeTrack(trackId); // API call to like the track
    
            } catch (error) {
                console.error('Error liking track:', error);
            }
        },
        async unlikeTrackHandler(trackId) {
            try {
                // Unlike the track locally and on the server
                this.likedTracks.delete(trackId); // Correct method to remove the track from the set
                this.$toast.add({
                    severity: 'info',
                    summary: 'Success',
                    detail: 'Track unliked successfully!',
                });
                await trackService.unlikeTrack(trackId); // API call to unlike the track
            } catch (error) {
                console.error('Error unliking track:', error);
            }
        },
        isLiked(trackId) {
            return this.likedTracks.has(trackId);
        },
    },
    mounted() {
        this.fetchRecommendedTracks();
        this.likedTracksHandler();
    },
};
</script>

<template>
    <div class="card">
        <DataTable :value="tracks" :rows="5" :paginator="true" responsiveLayout="scroll">
            <Column field="title" header="Title" style="width: 50%">
                <template #body="{ data }">
                    <div class="flex items-center gap-2">
                        <img :alt="data.track_image_url" :src="data.track_image_url" style="width: 32px" />
                        <span>{{ data.track_title }}</span>
                    </div>
                </template>
            </Column>
            <Column field="artist" header="Artist" style="width: 40%">
                <template #body="{ data }">
                    <div class="flex items-center gap-2">
                        <span>{{ data.artist_name }}</span>
                    </div>
                </template>
            </Column>
            <Column style="width: 10%" header="">
                <template #body="{ data }">
                    <div class="flex gap-2">
                        <!-- Play Button -->
                        <Button icon="pi pi-play" type="button" class="p-button-text"
                            @click="handleClick(data.track_uri)"></Button>

                        <!-- Like/Unlike Button -->
                        <Button 
                            :icon="isLiked(data.track_id) ? 'pi pi-heart-fill' : 'pi pi-heart'" 
                            type="button" 
                            class="p-button-text" 
                            @click="isLiked(data.track_id) ? unlikeTrackHandler(data.track_id) : likeTrackHandler(data.track_id)">
                        </Button>
                    </div>
                </template>
            </Column>
        </DataTable>
    </div>
</template>

<style scoped>
/* Optional: Add some custom styles for the like button if needed */
.pi-heart-fill {
    color: red;
}
</style>
