<template>
    <div class="update-tracks-dialog">
      <Dialog header="Update Tracks" :visible="visible" @hide="closeDialog">
        <div>
          <h5>Select Tracks to Add:</h5>
          <div v-if="tracks.length">
            <div v-for="track in tracks" :key="track.track_id" class="track-item">
              <Checkbox
                v-model="selectedTracks"
                :value="track"
                :label="track.track_name"
              />
            </div>
          </div>
          <div v-else>
            <p>No tracks available.</p>
          </div>
        </div>
  
        <template #footer>
          <Button label="Cancel" icon="pi pi-times" @click="closeDialog" />
          <Button label="Save" icon="pi pi-check" @click="saveTracks" :disabled="selectedTracks.length === 0" />
        </template>
      </Dialog>
    </div>
  </template>
  
  <script>

  import { API_BASE_URL } from "@/service/apiConfig";
import axios from "axios";
  
  export default {
    components: {

    },
    props: {
      visible: {
        type: Boolean,
        required: true,
      },
      playlistId: {
        type: String,
        required: true,
      },
      currentTracks: {
        type: Array,
        default: () => [],
      },
    },
    data() {
      return {
        tracks: [],
        selectedTracks: [],
        API_BASE_URL,
      };
    },
    methods: {
      async fetchTracks() {
        try {
          const response = await axios.get(`${this.API_BASE_URL}/tracks`, {
            withCredentials: true,
          });
          this.tracks = response.data;
        } catch (error) {
          console.error("Error fetching tracks:", error);
        }
      },
  
      closeDialog() {
        this.$emit("update:visible", false);
      },
  
      async saveTracks() {
        try {
          const trackIds = this.selectedTracks.map(track => track.track_id);
          await axios.post(`${this.API_BASE_URL}/playlists/${this.playlistId}/update-tracks`, { trackIds }, {
            withCredentials: true,
          });
  
          // Emit an event to notify the parent component to update the playlist
          this.$emit("tracks-updated", trackIds);
          this.closeDialog();
        } catch (error) {
          console.error("Error saving tracks:", error);
        }
      },
    },
    watch: {
      visible(newVal) {
        if (newVal) {
          this.fetchTracks();
        }
      },
    },
  };
  </script>
  
  <style scoped>
  .update-tracks-dialog {
    width: 500px;
  }
  .track-item {
    margin-bottom: 10px;
  }
  </style>
  