<template>
  <div>
    <!-- Loading Spinner -->
    <div
      ref="loadingSpinner"
      v-if="loading"
      class="p-d-flex p-jc-center p-ai-center"
    >
      <span>Loading...</span>
    </div>

    <!-- Content Wrapper: Events and Playlists -->
    <div
      ref="contentWrapper"
      v-if="!loading && (user.events.length || user.playlists.length)"
    >
      <!-- Display Events -->
      <div class="p-4" v-if="user.events.length">
        <h2 class="text-xl font-semibold mb-3">Events</h2>
        <div class="flex space-x-4 overflow-x-auto pb-4 h-full">
          <EventComponent
            v-for="event in user.events"
            :key="event.event_id"
            :event="event"
          />
        </div>
      </div>

      <!-- Display Playlists -->
      <div class="p-4" v-if="user.playlists.length">
        <h2 class="text-xl font-semibold mb-3">Playlists</h2>
        <div class="flex space-x-4 overflow-x-auto pb-4">
          <PlaylistComponent
            v-for="playlist in user.playlists"
            :key="playlist.playlist_id"
            :playlist="playlist"
          />
        </div>
      </div>
    </div>

    <!-- No Events or Playlists Found -->
    <div
      v-if="!loading && !user.events.length && !user.playlists.length"
      class="p-4"
    >
      <p>No events or playlists found.</p>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import { API_BASE_URL } from "@/service/apiConfig";
import EventComponent from "@/components/EventComponent.vue";
import PlaylistComponent from "@/components/PlaylistComponent.vue";

export default {
  components: {
    EventComponent,
    PlaylistComponent,
  },
  data() {
    return {
      loading: false,
      user: { events: [], playlists: [] },
      errorMessage: "",
    };
  },
  methods: {
    async getEventsByUserId() {
      try {
        const response = await axios.get(`${API_BASE_URL}/me/events`, {
          withCredentials: true,
        });
        console.log("Fetched events:", response.data);
        if (Array.isArray(response.data)) {
          this.user.events = response.data;
        } else {
          throw new Error("Invalid events data format");
        }
      } catch (err) {
        this.handleError(err, "events");
      }
    },

    async getPlaylistsForUser() {
      try {
        const response = await axios.get(`${API_BASE_URL}/me/playlists`, {
          withCredentials: true,
        });
        console.log("Fetched playlists:", response.data);
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
      this.errorMessage =
        error.response?.data?.message || `Failed to fetch ${dataType}.`;
    },

    async fetchEvents() {
      this.errorMessage = "";
      this.loading = true;
      console.log("Fetching events and playlists...");
      try {
        await Promise.all([
          this.getEventsByUserId(),
          this.getPlaylistsForUser(),
        ]);
        console.log("Final user object:", this.user);
      } catch (err) {
        // No need to handle here, errors handled in individual methods
      } finally {
        this.loading = false;
      }
    },
  },
  mounted() {
    this.fetchEvents();
  },
};
</script>
