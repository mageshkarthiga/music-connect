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
        <div class="flex space-x-4 overflow-x-auto pb-4">
          <EventComponent v-for="event in user.events" :key="event.event_id" :event="event" />
        </div>
      </div>

      <!-- Display Playlists -->
      <div class="p-4" v-if="user.playlists.length">
        <h2 class="text-xl font-semibold mb-3">Playlists</h2>
        <div class="flex space-x-4 overflow-x-auto pb-4">
          <PlaylistComponent v-for="playlist in user.playlists" :key="playlist.playlist_id" :playlist="playlist" />
        </div>
      </div>

      <div class="p-4" v-if="user.playlists.length">
          <RecommendedTracks />
      </div>
    </div>

    <!-- No Events or Playlists Found -->
    <div v-if="!loading && !user.events.length && !user.playlists.length" class="p-4">
      <p>No events or playlists found.</p>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import { API_BASE_URL } from "@/service/apiConfig";
import EventComponent from "@/components/EventComponent.vue"; // Import EventComponent
import PlaylistComponent from "@/components/PlaylistComponent.vue"; 
import RecommendedTracks from "@/components/RecommendedTracks.vue";

export default {
  components: {
    EventComponent, // Register EventComponent
    PlaylistComponent, // Register PlaylistComponent
    RecommendedTracks, 
  },
  data() {
    return {
      loading: false, // Flag to manage loading state
      user: { events: [], playlists: [] }, // Store user data
      errorMessage: "", // Store error messages
    };
  },
  methods: {
    // Fetch events by user ID with enhanced error handling
    async getEventsByUserId() {
      try {
        const response = await axios.get(`${API_BASE_URL}/me/events`, {
          withCredentials: true,
        });
        if (Array.isArray(response.data)) {
          this.user.events = response.data;
        } else {
          throw new Error("Invalid events data format");
        }
      } catch (err) {
        this.handleError(err, "events");
      }
    },

    // Fetch playlists for user with enhanced error handling
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

    // Error handling
    handleError(error, dataType) {
      console.error(`${dataType} fetch error:`, error);
      this.errorMessage = error.response?.data?.message || `Failed to fetch ${dataType}.`;
    },

    // Fetch both events and playlists
    async fetchEvents() {
      this.errorMessage = "";
      this.loading = true;
      try {
        await Promise.all([
          this.getEventsByUserId(),
          this.getPlaylistsForUser(),
        ]);
      } catch (err) {
        // Handle any fetch error
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
