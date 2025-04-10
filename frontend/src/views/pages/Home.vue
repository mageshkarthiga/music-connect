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

    <!-- Error Message -->
    <div v-if="errorMessage" class="p-4 text-red-500">
      {{ errorMessage }}
    </div>

    <!-- Content -->
    <div v-if="!loading">
      <template v-if="hasContent">
        <!-- Events -->
        <div class="p-4" v-if="Array.isArray(user.events) && user.events.length">
          <h2 class="text-xl font-semibold mb-3">Events</h2>
          <div class="flex space-x-4 overflow-x-auto pb-4 h-full">
            <EventCard
              v-for="event in user.events"
              :key="event.event_id"
              :event="event"
            />
          </div>
        </div>

        <!-- Playlists -->
        <div class="p-4" v-if="Array.isArray(user.playlists) && user.playlists.length">
          <h2 class="text-xl font-semibold mb-3">Playlists</h2>
          <div class="flex space-x-4 overflow-x-auto pb-4">
            <PlaylistCard
              v-for="playlist in user.playlists"
              :key="playlist.playlist_id"
              :playlist="playlist"
            />
          </div>
        </div>

        <!-- Tracks -->
        <div class="p-4" v-if="Array.isArray(user.tracks) && user.tracks.length">
          <h2 class="text-xl font-semibold mb-3">Tracks</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            <TrackCard
              v-for="track in user.tracks"
              :key="track.track_id"
              :track="track"
            />
          </div>
        </div>

        <!-- Artists -->
        <div class="p-4" v-if="user.tracks.length">
          <h2 class="text-xl font-semibold mb-3">Artists</h2>
          <SpotifyPlayer />
        </div>

        <div class="p-4">
          <div class="font-semibold text-xl mb-4">Recommended music</div>
          <RecommendedTracks />
        </div>
      </template>

      <!-- No Content -->
      <template v-else>
        <div class="p-4">
          <p>No events, playlists, or tracks found.</p>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import { API_BASE_URL } from "@/service/apiConfig";
import EventCard from "@/components/EventCard.vue";
import PlaylistCard from "@/components/PlaylistCard.vue";
import TrackCard from "@/components/TrackCard.vue";
import SpotifyPlayer from "@/components/SpotifyPlayer.vue";
import RecommendedTracks from "@/components/RecommendedTracks.vue";

export default {
  components: {
    EventCard,
    PlaylistCard,
    TrackCard,
    SpotifyPlayer,
    RecommendedTracks
  },
  data() {
    return {
      loading: false,
      user: { events: [], playlists: [], tracks: [] },
      errorMessage: "",
    };
  },
  computed: {
    hasContent() {
      return (
        this.user.events.length ||
        this.user.playlists.length ||
        this.user.tracks.length
      );
    },
  },
  methods: {
    async getEventsByUserId() {
      try {
        const response = await axios.get(`${API_BASE_URL}/me/events`, {
          withCredentials: true,
        });
        this.user.events = response.data;
      } catch (err) {
        this.handleError(err, "events");
      }
    },
    async getPlaylistsForUser() {
      try {
        const response = await axios.get(`${API_BASE_URL}/me/playlists`, {
          withCredentials: true,
        });
        this.user.playlists = response.data;
      } catch (err) {
        this.handleError(err, "playlists");
      }
    },
    async getTracksForUser() {
      try {
        const response = await axios.get(`${API_BASE_URL}/me/tracks`, {
          withCredentials: true,
        });
        this.user.tracks = response.data;
      } catch (err) {
        this.handleError(err, "tracks");
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
      try {
        await Promise.all([
          this.getEventsByUserId(),
          this.getPlaylistsForUser(),
          this.getTracksForUser(),
        ]);
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
