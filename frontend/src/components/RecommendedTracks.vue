<script setup>
import { ProductService } from '@/service/ProductService';
import { getRecommendedTracks } from '@/service/RecommenderService';
import { onMounted, ref } from 'vue';

const tracks = ref(null);

onMounted(() => {
    getRecommendedTracks().then((data) => {
        tracks.value = data;
    });
});
</script>

<template>
    <div class="card">
        <div class="font-semibold text-xl mb-4">Recommended music</div>
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
                <template #body>
                    <Button icon="pi pi-caret-right" type="button" class="p-button-text"></Button>
                </template>
            </Column>
        </DataTable>
    </div>
</template>
