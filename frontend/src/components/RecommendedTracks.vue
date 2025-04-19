<script>
import { getRecommendedTracks } from '@/service/RecommenderService';

export default {
    name: 'RecommendedTracks',
    data() {
        return {
            tracks: null,
            selectedTrackUri: null,
        };
    },
    methods: {
        async fetchRecommendedTracks() {
            try {
                const data = await getRecommendedTracks();
                this.tracks = data;
            } catch (error) {
                console.error('Error fetching recommended tracks:', error);
            }
        },
        handleClick(trackUri) {
            this.selectedTrackUri = trackUri;
            this.$emit('track-selected', trackUri);
        },
    },
    mounted() {
        this.fetchRecommendedTracks();
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
                    <Button icon="pi pi-play" type="button" class="p-button-text"
                        @click="handleClick(data.track_uri)"></Button>
                </template>
            </Column>
        </DataTable>
    </div>
</template>