<template>
  <!-- Sidebar -->
  <div class="sidebar">
    <!-- Library Header -->
    <div class="library-header">
      <h4>Your Library</h4>
      <Button @click="openAddPlaylistDialog" class="add-playlist-btn">
        <i class="pi pi-plus"></i>
      </Button>
    </div>

    <!-- Playlists List -->
    <div v-if="user.playlists.length" class="playlist-list">
      <div
        v-for="playlist in user.playlists"
        :key="playlist.playlist_id"
        class="playlist-item"
        @click="navigateToPlaylist(playlist.playlist_id, playlist.playlist_name)"
      >
        <Button
          @click="openPlaylist(playlist.playlist_id)"
          class="playlist-button"
          label="Secondary"
          severity="secondary"
          text
        >
          <img :src="playlist.playlist_image_url || `https://i.scdn.co/image/ab67616d0000b2735074bd0894cb1340b8d8a678`" alt="Playlist Image" class="playlist-image" />
          <div class="playlist-info">
            <span class="playlist-name">{{ playlist.playlist_name }}</span>
            <span class="playlist-username"> Playlist - {{ currentUser?.user_name }}</span>
          </div>
        </Button>
      </div>
    </div>

    <div v-else>
      <p>No playlists available.</p>
    </div>

    <!-- Playlist Dialog -->
    <AddPlaylistDialog
      v-if="showPlaylistDialog"
      :showDialog="showPlaylistDialog"
      @update:showDialog="showPlaylistDialog = $event"
      @save="handleSave"
      @close="closeAddPlaylistDialog"
    />
  </div>
</template>

<script>
import axios from "axios";
import AddPlaylistDialog from "@/components/AddPlaylistDialog.vue";

export default {
  components: {
    AddPlaylistDialog,
  },
  data() {
    return {
      user: { playlists: [] },
      currentUser: null,
      errorMessage: "",
      darkTheme: localStorage.getItem("theme") === "dark",
      showPlaylistDialog: false,
      API_BASE_URL: process.env.VUE_APP_API_BASE_URL,
    };
  },
  methods: {
    async getCurrentUser() {
      try {
        const response = await axios.get(`${API_BASE_URL}/me`, {
          withCredentials: true,
        });
        this.currentUser = response.data;
      } catch (err) {
        this.handleError(err, "current user");
      }
    },

    async getPlaylistsForUser() {
      try {
        const response = await axios.get(`${API_BASE_URL}/me/playlists`, {
          withCredentials: true,
        });

        if (Array.isArray(response.data)) {
          this.user.playlists = response.data;
        } else {
          throw new Error("Invalid playlists data format");
        }
      } catch (err) {
        this.handleError(err, "playlists");
      }
    },

    handleError(error, dataType) {
      console.error(`${dataType} fetch error:`, error);
      this.errorMessage = error.response?.data?.message || `Failed to fetch ${dataType}.`;
    },

    openAddPlaylistDialog() {
      this.showPlaylistDialog = true;
    },

    closeAddPlaylistDialog() {
      this.showPlaylistDialog = false;
    },

    handleSave(newPlaylist) {
      this.user.playlists.push(newPlaylist);
      this.closeAddPlaylistDialog();
    },

    openPlaylist(playlistId) {
      this.$router.push({ name: "Playlist", query: { id: playlistId } });
    },

    navigateToPlaylist(playlist_id, playlist_name) {
      this.$router.push({
        name: "playlist",
        params: {
          playlist_id: playlist_id,
          playlist_name: playlist_name,
        },
      });
    },
  },
  mounted() {
    this.getCurrentUser();
    this.getPlaylistsForUser();
  },
};
</script>

<style lang="scss" scoped>
.sidebar {
  padding: 1rem;
  width: 250px;
  height: auto;
  max-height: 92.5%;
  display: flex;
  flex-direction: column;
  overflow-y: auto;
  transition: background-color 0.3s ease, color 0.3s ease;
}

.dark {
  background-color: #121212;
  color: white;
}

.light {
  background-color: #f4f4f4;
  color: black;
}

h4 {
  color: inherit;
}

.library-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.add-playlist-btn {
  background-color: #1db954;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 50%;
  cursor: pointer;
  font-size: 1.5rem;
  width: 40px;
  height: 40px;
}

.add-playlist-btn:hover {
  background-color: #1ed760;
}

.theme-toggle-btn {
  background-color: transparent;
  border: none;
  color: inherit;
  font-size: 1.5rem;
  cursor: pointer;
}

.playlist-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.playlist-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem;
  border-radius: 8px;
  transition: background-color 0.3s ease;
}

.dark .playlist-item {
  background-color: #282828;
}

.light .playlist-item {
  background-color: #ffffff;
}

.playlist-item:hover {
  background-color: rgba(125, 125, 125, 0.1);
}

.dark .playlist-item:hover {
  background-color: #3a3a3a;
}

.light .playlist-item:hover {
  background-color: #cacaca;
}

.playlist-image {
  width: 50px;
  height: 50px;
  border-radius: 8px;
  object-fit: cover;
}

.playlist-info {
  display: flex;
  flex-direction: column;
}

.playlist-name {
  font-size: 1rem;
  font-weight: 600;
}

.playlist-username {
  font-size: 0.875rem;
  color: rgba(125, 125, 125, 0.9);
}

.dark .playlist-name {
  color: rgba(255, 255, 255, 0.7);
}

.light .playlist-username {
  color: rgb(0, 0, 0, 0.7);
}
</style>
