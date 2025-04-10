<template>
  <div class="dialog-overlay">
    <div class="dialog-box">
      <h3>Create Playlist</h3>

      <!-- Playlist Name -->
      <div>
        <label for="playlistName">Playlist Name</label>
        <input
          v-model="playlistName"
          id="playlistName"
          type="text"
          placeholder="Enter playlist name"
        />
      </div>

      <!-- Track Selection -->
      <div class="track-list">
        <TrackCard
          v-for="track in tracks"
          :key="track.track_id"
          :track="track"
          :state="'select'"  
          :selectedTracks="selectedTracks"
          @toggle="toggleTrack"
        />
      </div>

      <!-- Buttons -->
      <Button @click="savePlaylist" :disabled="!playlistName || !selectedTracks.length">
        Save Playlist
      </Button>
      <Button @click="closeDialog">Cancel</Button>
    </div>
  </div>
</template>

<script>
import PlaylistService from "@/service/PlaylistService"; // Import PlaylistService
import TrackCard from "@/components/TrackCard.vue"; // Import the TrackCard component
import userService from "@/service/UserService";
import axios from "axios"; // Import axios for fetching tracks
import { API_BASE_URL } from "@/service/apiConfig"; // Import API base URL

export default {
  components: {
    TrackCard, // Register the TrackCard component
  },
  props: {
    showDialog: Boolean,
    onClose: Function,
    onSave: Function,
  },
  data() {
    return {
      playlistName: "",
      selectedTracks: [], // Array to store selected track IDs
      tracks: [], // Array to hold all the fetched tracks
      userId: null, // Store userId here
    };
  },
  methods: {
    async getUser() {
      try {
        const user = await userService.getUser();
        this.userId = user.user_id;
      } catch (error) {
        console.error("Error fetching user:", error);
      }
    },

    async fetchTracks() {
      try {
        const response = await axios.get(`${API_BASE_URL}/tracks`, {
          withCredentials: true,
        });
        this.tracks = response.data; // Store fetched tracks in tracks array
      } catch (error) {
        console.error("Error fetching tracks:", error);
      }
    },

    // This method will be used to toggle track selection
    toggleTrack(trackId) {

      const index = this.selectedTracks.indexOf(trackId);
      console.log("Track ID:", trackId);
      if (index > -1) {
        // Remove from selected
        this.selectedTracks.splice(index, 1);
      } else {
        // Add to selected
        this.selectedTracks.push(trackId);
      }
    },

    async savePlaylist() {
 

      const trackIds = this.selectedTracks;
      console.log("Selected track IDs:", trackIds);
      try {
        const newPlaylist = await PlaylistService.addPlaylistForUser(this.userId, this.playlistName, trackIds);
        if (this.onSave) {
          this.onSave(newPlaylist); // Pass new playlist back to parent component
        }
        this.closeDialog();
      } catch (error) {
        console.error("Error saving playlist:", error);
      }
    },

    closeDialog() {
      if (this.onClose) {
        this.onClose(); // Close the dialog when done
      }
    },
  },
  mounted() {
    this.getUser(); // Fetch user details when the dialog is mounted
    this.fetchTracks(); // Fetch all tracks when the dialog is mounted
  },
};
</script>

<style scoped>
.dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.dialog-box {
  background-color: white;
  padding: 20px;
  border-radius: 8px;
  width: 600px;
  max-height: 60vh;
  overflow-y: auto;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
}

.track-list {
  margin-top: 10px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
}
</style>
