<template>
  <div :class="'dialog-overlay'">
    <div :class="['dialog-box', { 'dark-theme': isDarkMode }]">
      <h3 :class="{'light': !isDarkMode, 'dark': isDarkMode}">Create Playlist</h3>

      <!-- Playlist Name -->
      <div>
        <label for="playlistName">Playlist Name:  </label>
        <input
          v-model="playlistName"
          id="playlistName"
          type="text"
          placeholder="Enter playlist name"
          :class="{'light': !isDarkMode, 'dark': isDarkMode}"
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
      <Button @click="savePlaylist" :disabled="!playlistName || !selectedTracks.length ">
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
    isDarkMode: Boolean, // New prop to track theme mode
  },
  data() {
    return {
      playlistName: "",
      selectedTracks: [], // Array to store selected track IDs
      tracks: [], // Array to hold all the fetched tracks
      userId: null, // Store userId here
      trackImageUrl: "", // Store track image URL here
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
      console.log("Selected Track ID:", trackId);

      if (index > -1) {
        // Remove from selected
        this.selectedTracks.splice(index, 1);
      } else {
        // Add to selected
        this.selectedTracks.push(trackId);
      }

      console.log("Selected Tracks Array:", this.selectedTracks);

      // Set the track image URL to the first selected track's image URL (if there are selected tracks)
      if (this.selectedTracks.length > 0) {
        const firstSelectedTrack = this.tracks.find(track => track.track_id === this.selectedTracks[0]);
        if (firstSelectedTrack) {
          this.trackImageUrl = firstSelectedTrack.track_image_url; // Assuming the image URL is stored in `track_image_url`
        }
      } else {
        this.trackImageUrl = ""; // Clear the image URL if no tracks are selected
      }
    },

    async savePlaylist() {
  // Convert the hash map to an array of selected track IDs
  const trackIds = Object.keys(this.selectedTracks)
    .filter(trackId => this.selectedTracks[trackId])
    .map(id => parseInt(id)); // ensure IDs are numbers

  console.log("Selected track IDs:", trackIds);

  let newPlaylist;

  try {
    // Step 1: Create the playlist
    newPlaylist = await PlaylistService.createPlaylistForUser(
      this.playlistName,        // name
      this.trackImageUrl,       // playlistImageUrl
      this.userId               // userId
    );

    console.log("New Playlist created:", newPlaylist);

  } catch (error) {
    console.error("Error creating playlist:", error);
    return;
  }

  try {
    const playlistId = newPlaylist.playlist_id; // Now it's defined properly
    console.log("New Playlist ID:", playlistId);
    console.log("Track IDs to add:", trackIds);
    
    // Step 2: Add tracks to the newly created playlist
    if (trackIds.length > 0) {
      await PlaylistService.addTracksToPlaylist(playlistId, trackIds);
      console.log("Tracks successfully added to playlist.");
    }

    // Step 3: Notify parent and close dialog
    if (this.onSave) {
      this.onSave(newPlaylist); // Notify parent
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

.track-list {
  margin: 200px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr)); /* Increased minimum width for better spacing */
  gap: 24px; /* Increased the gap between track cards */
}

.track-card {
  border-radius: 12px;
  cursor: pointer;
  transition: 0.2s;
}

.dark .track-list{
  background-color: #2d3748; /* Darker background */
  color: #fff; /* White text */
  margin: 0;

}

.light .track-list{
  background-color: #f8f9fa; /* Light background */
  color: #000; /* Black text */
}

.dark .dialog-box {
  background-color: #2d3748; /* Darker background */
  color: #fff; /* White text */
}

.light .dialog-box {
  background-color: #f8f9fa; /* Light background */
  color: #000; /* Black text */
}

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

/* Light mode styles */
.dialog-box {
  background-color: white;
  color: black;
  padding: 20px;
  border-radius: 8px;
  width: 600px;
  max-height: 60vh;
  overflow-y: auto;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
}

h3.light {
  color: black;
  margin-bottom: 20px;
}

h3.dark {
  color: white;
}

.light .dialog-box {
  background-color: #f8f9fa;
  color: #212529;
}

.dark .dialog-box {
  background-color: #343a40;
  color: white;
}

.track-list {
  margin-top: 10px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
}

/* Dark mode styles */
.dialog-overlay.dark-theme {
  background-color: rgba(0, 0, 0, 0.7);
}

.dialog-box.dark-theme {
  background-color: #333;
  color: white;
}

.dialog-box.dark-theme button {
  background-color: #444;
  color: white;
  border: 1px solid #666;
}

.dialog-box.dark-theme input {
  background-color: #555;
  color: white;
  border: 1px solid #777;
}

/* Input field styles for light and dark themes */
input.light {
  background-color: #fff;
  color: #000;
  border: 1px solid #ccc;
}

input.dark {
  background-color: #444;
  color: #fff;
  border: 1px solid #666;
}

input::placeholder {
  color: #aaa;
}

/* Placeholder for light and dark modes */
input.light::placeholder {
  color: #aaa;
}

input.dark::placeholder {
  color: #888;
}

/* Dark mode styles */
.dialog-overlay.dark-theme {
  background-color: rgba(0, 0, 0, 0.7);
}

.dialog-box.dark-theme {
  background-color: #2d3748; /* Darker background for the dialog box */
  color: white; /* White text for better contrast */
}

h3.dark {
  color: white; /* White text for dark mode */
}

.light .dialog-box {
  background-color: #f8f9fa;
  color: #212529;
}

.dark .dialog-box {
  background-color: #343a40; /* Dark background for dialog box */
  color: white; /* Light text in dark mode */
}

.track-list {
  margin-top: 10px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
}

/* Dark mode track-list styles */
.track-list.dark {
  background-color: #2d3748; /* Darker background for track-list */
  color: white; /* Light text for track list items */
}

/* Input field styles for light and dark themes */
input.light {
  background-color: #fff;
  color: #000;
  border: 1px solid #ccc;
}

input.dark {
  background-color: #444;
  color: #fff;
  border: 1px solid #666;
}

/* Input placeholder styles for light and dark modes */
input.light::placeholder {
  color: #aaa;
}

input.dark::placeholder {
  color: #888; /* Slightly lighter placeholder text for dark mode */
}

button.dark {
  background-color: #444;
  color: white;
  border: 1px solid #666;
}

button.light {
  background-color: #f8f9fa;
  color: #212529;
  border: 1px solid #ccc;
}

.dark .dialog-box {
  background-color: #2d3748; /* Darker background */
  color: white; /* White text */
}

.light .dialog-box {
  background-color: #f8f9fa; /* Light background */
  color: #212529; /* Dark text */
}


.dialog-box.dark-theme button {
  background-color: #444;
  color: white;
  border: 1px solid #666;
}

.dialog-box.dark-theme input {
  background-color: #555;
  color: white;
  border: 1px solid #777;
}

.dark  .dialog-overlay {
  background-color: rgba(0, 0, 0, 0.7);

}

.light .dialog-overlay {
  background-color: rgba(255, 255, 255, 0.7);
}

/* Global styles for dialog box */
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
  padding: 20px;
  border-radius: 8px;
  width: 700px;
  max-height: 60vh;
  overflow-y: auto;
  box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
  color: black; /* Default text color */
}

h3 {
  margin-bottom: 20px;
}

.track-list {
  margin-top: 10px;
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
}

/* Light mode styles */
.dialog-overlay.light-theme {
  background-color: rgba(255, 255, 255, 0.7);
}

.dialog-box.light-theme {
  background-color: #f8f9fa;
  color: #212529;
}

.light input {
  background-color: #fff;
  color: #000;
  border: 1px solid #ccc;
}

.light button {
  background-color: #f8f9fa;
  color: #212529;
  border: 1px solid #ccc;
}

/* Dark mode styles */
.dialog-overlay.dark-theme {
  background-color: rgba(0, 0, 0, 0.7);
}

.dialog-box.dark-theme {
  background-color: #2d3748;
  color: white;
}

.dark  h3. {
  color: white;
}

.dark  input{
  background-color: #444;
  color: #fff;
  border: 1px solid #666;
}

.dark .button {
  background-color: #444;
  color: white;
  border: 1px solid #666;
}

.dark .track-list {
  background-color: #2d3748;
  color: white;
}

.light .track-list {
  background-color: #f8f9fa;
  color: #212529;
}


</style>
