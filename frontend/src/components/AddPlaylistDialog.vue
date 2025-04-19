<template>
  <Dialog
    :visible="showDialog"
    :modal="true"
    :closable="true"
    :breakpoints="{'960px': '75vw', '640px': '90vw'}"
    :style="{ width: '50vw' }"
    header="Create Playlist"
    @hide="closeDialog"
  >
    <!-- Playlist Name -->
    <div class="p-fluid">
      <div class="field">
        <label for="playlistName">Playlist Name: </label>
        <InputText
          v-model="playlistName"
          id="playlistName"
          placeholder="Enter playlist name"
        />
      </div>
    </div>
 

    <!-- Track Selection -->
    <h6>Select Tracks: </h6>
    <div class="track-list">

      <div v-if="tracks.length > 0">
        <TrackCard
          v-for="track in tracks"
          :key="track.track_id"
          :track="track"
          :state="'select'"
          :selectedTracks="selectedTracks"
          @toggle="toggleTrack"
        />
      </div>
      <div v-else>
        <p>No tracks available to display.</p>
      </div>
      <div v-else>
        <p>No tracks available to display.</p>
      </div>
    </div>


    <!-- Buttons -->
    <template #footer>
      <Button
        label="Save Playlist"
        icon="pi pi-check"
        @click="savePlaylist"
        :disabled="!playlistName || !selectedTracks.length"
        class="p-button-success"
      />
      <Button
        label="Cancel"
        icon="pi pi-times"
        @click="closeDialog"
        class="p-button-secondary"
      />
    </template>
  </Dialog>
</template>

<script>
import PlaylistService from "@/service/PlaylistService";
import TrackCard from "@/components/TrackCard.vue";
import userService from "@/service/UserService";
import trackService from "@/service/TrackService";

export default {
  components: {
    TrackCard,
  },
  props: {
    showDialog: Boolean,
    onSave: Function,
    onSave: Function,
    onClose: Function,
    playlistToEdit: Object,
    playlistToEdit: Object,
  },
  data() {
    return {
      playlistName: "",
      selectedTracks: [],
      tracks: [],
      userId: null,
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
        const tracks = await trackService.getTracks();
        this.tracks = tracks;
      } catch (error) {
        console.error("Error fetching tracks:", error);
      }
    },

    toggleTrack(trackId) {
      const index = this.selectedTracks.indexOf(trackId);

      if (index > -1) {
        this.selectedTracks.splice(index, 1);
      } else {
        this.selectedTracks.push(trackId);
      }

      // if (this.selectedTracks.length > 0) {
      //   const firstSelectedTrack = this.tracks.find(track => track.track_id === this.selectedTracks[0]);
      //   if (firstSelectedTrack) {
      //     this.trackImageUrl = firstSelectedTrack.track_image_url;
      //   }
      // } else {
      //   this.trackImageUrl = "";
      // }
    },

    async savePlaylist() {
      if (!this.userId) {
        console.error("User ID is missing");
        return;
      }

      let newPlaylist;

      try {
        newPlaylist = await PlaylistService.createPlaylistForUser(
          this.playlistName,
          this.trackImageUrl,
          this.userId
        );
      } catch (error) {
        console.error("Error creating playlist:", error);
        return;
      }

      try {
        const playlistId = newPlaylist.playlist_id;
        if (this.selectedTracks.length > 0) {
          await PlaylistService.addTracksToPlaylist(playlistId, this.selectedTracks);
        }

        if (this.onSave) {
          this.onSave(newPlaylist);
        }

        this.$emit("playlist-added", newPlaylist);
        this.closeDialog();
      } catch (error) {
        console.error("Error saving playlist:", error);
      }
    },

    closeDialog() {
      if (this.onClose) {
        this.onClose();
        this.onClose();
      }
    },
  },
  mounted() {
    this.getUser();
    this.fetchTracks();
  },

  computed: {
  trackImageUrl() {
    if (this.selectedTracks.length === 0) return "";
    const first = this.tracks.find(t => t.track_id === this.selectedTracks[0]);
    return first?.track_image_url || "";
  },
},


  watch: {
    playlistToEdit: {
      immediate: true,
      handler(playlist) {
        if (playlist) {
          this.playlistName = playlist.name;
          // this.trackImageUrl = playlist.image_url;
          this.selectedTracks = playlist.tracks.map(t => t.track_id);
        }
      }
    }
  },
};


</script>

<style scoped>
.track-list {
  color: black;
  display: flex;
  flex-wrap: wrap;
  gap: 30px;
  padding: 10px;
  margin: 10px;
  justify-content: center; /* Centers the tracks horizontally */
  align-items: center; /* Centers the tracks vertically */
}
</style>
