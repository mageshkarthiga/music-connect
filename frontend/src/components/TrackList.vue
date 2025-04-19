<template>
    <div class="track-list">
      <h2>Tracks</h2>
      <TrackCard
          v-for="track in tracks"
          :key="track.track_id"
          :track="track"
          :state="'select'"  
          :selectedTracks="selectedTracks"
          @toggle="toggleTrack"
        />
    </div>
  </template>
  
  <script>
  import TrackCard from "./TrackCard.vue";
  
  export default {
    name: "TrackList",
    components: { TrackCard },
    props: {
      tracks: Array,
    },
  };
  </script>
  
  <style scoped>
  .track-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  </style>
  