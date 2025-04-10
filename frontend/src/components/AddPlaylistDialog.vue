<template>
  <div class="dialog-overlay">
    <div class="dialog-box">
      <h3>Create Playlist</h3>
      
      <!-- Playlist Name -->
      <div>
        <label for="playlistName">Playlist Name</label>
        <input v-model="playlistName" id="playlistName" type="text" placeholder="Enter playlist name" />
      </div>

      <!-- Track Selection -->
      <track-list :selectedTracks="selectedTracks" />

      <button @click="savePlaylist" :disabled="!playlistName || !selectedTracks.length">
        Save Playlist
      </button>
      <button @click="closeDialog">Cancel</button>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import TrackList from "@/components/TrackList.vue"; // Import the TrackList component
import userService from "@/service/UserService";

export default {
  components: {
    TrackList, // Register the TrackList component
  },
  props: {
    showDialog: Boolean,
    onClose: Function,
    onSave: Function,
  },
  data() {
    return {
      playlistName: "",
      username: "",
      selectedTracks: [],
    };
  },
  methods: {
    async getUser() {
      try {
        // Using the getUser method from userService
        const user = await userService.getUser();
        this.username = user.username; // Set the username from the userService response
      } catch (error) {
        console.error("Error fetching user:", error);
      }
    },

    async savePlaylist() {
      const playlistData = {
        name: this.playlistName,
        username: this.username,
        tracks: this.selectedTracks,
      };

      if (this.onSave) {
        this.onSave(playlistData);
      }
    },

    closeDialog() {
      if (this.onClose) {
        this.onClose();
      }
    },
  },
  mounted() {
    this.getUser(); // Fetch user details when the dialog is mounted
  },
};
</script>
