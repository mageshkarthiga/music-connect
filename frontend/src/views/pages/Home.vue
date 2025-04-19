<template>
  <div>
    <!-- Loading Spinner -->
    <div ref="loadingSpinner" v-if="loading" class="p-d-flex p-jc-center p-ai-center">
      <span>Loading...</span>
    </div>

    <!-- Error Message -->
    <div v-if="errorMessage" class="p-4 text-red-500">
      {{ errorMessage }}
    </div>

    <!-- Content -->
    <div v-if="!loading">
      <template v-if="hasContent">

        <!-- Tracks Section -->
        <div class="p-6 bg-gradient-to-r from-green-400 via-green-400 to-green-500 rounded-lg shadow-lg mb-10"
          v-if="user.tracks.length">
          <!-- Title with Emoji -->
          <h2 class="text-3xl font-extrabold text-center text-gray-800 mb-4">
            🏆 Top Played Tracks
          </h2>

          <!-- Subheading with Emoji -->
          <p class="text-center text-lg  text-gray-700 mb-6 opacity-90">
            🌟 Your most played tracks, right here!
          </p>

          <!-- Track Cards Grid -->
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 gap-8 auto-rows-fr">
            <!-- Track Card Loop -->
            <TrackCard v-for="track in user.tracks" :key="track.track_id" :track="track" :state="'redirect'"
              @track-selected="setSelectedTrackURI" @click="incrementPlayCount(track.track_id)"
              class="bg-white p-4 rounded-lg">
              <div class="flex items-center justify-between">
                <!-- Track Info -->
                <div class="flex items-center space-x-3">
                  <!-- Track Title -->
                  <h3 class="text-lg font-semibold text-gray-900">
                    {{ track.track_title }}
                  </h3>
                </div>
                <!-- Play Count Icon -->
                <div class="text-sm text-gray-600">
                  🌟 {{ track.play_count }}
                </div>
              </div>
            </TrackCard>
          </div>
        </div>

        <!-- No Tracks Message -->
        <div v-if="!user.tracks.length"
          class="p-6 bg-gradient-to-r from-green-400 via-green-400 to-green-500 rounded-lg shadow-xl text-center mb-10">
          <h2 class="text-3xl font-extrabold text-gray-800 mb-4">
            ❗️ No Top Tracks Yet
          </h2>
          <p class="text-lg font-medium text-gray-700 mb-4">
            Start listening to your favourite tunes and watch your top tracks appear right here! 🎧
          </p>
          <p class="text-sm text-gray-600 italic mb-6">
            Every track counts! Start playing your music and track your journey. 🎶
          </p>
        </div>



        <!-- Events -->
        <div class="p-4" v-if="Array.isArray(user.events) && user.events.length">
          <h2 class="text-xl font-semibold mb-3">Events</h2>
          <div class="flex space-x-4 overflow-x-auto pb-4 h-full">
            <EventCard v-for="event in user.events" :key="event.event_id" :event="event" />
          </div>
        </div>

        <!-- Playlists -->
        <div class="p-4" v-if="Array.isArray(user.playlists) && user.playlists.length">
          <h2 class="text-xl font-semibold mb-3">Playlists</h2>
          <div class="flex space-x-4 overflow-x-auto pb-4">
            <PlaylistCard v-for="playlist in user.playlists" :key="playlist.playlist_id" :playlist="playlist" />
          </div>
        </div>

        <!-- Artists -->
        <!-- <div class="p-4" v-if="user.tracks.length">
          <h2 class="text-xl font-semibold mb-3">Artists</h2>
        </div> -->

        <div class="p-4">
          <div class="font-semibold text-xl mb-4">Recommended music</div>
          <RecommendedTracks @track-selected="setSelectedTrackURI" />
        </div>
      </template>
    </div>

    <!-- No Content -->
    <template v-else>
      <div class="p-4">
        <p>No events, playlists, or tracks found.</p>
      </div>
    </template>

    <!-- Spotify Player -->
    <SpotifyPlayer v-if="selectedTrackURI" :spotifyUri="selectedTrackURI" />
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
import { incrementTrackPlayCount } from "@/service/TrackService";

export default {
  components: {
    EventCard,
    PlaylistCard,
    TrackCard,
    SpotifyPlayer,
    RecommendedTracks,
  },
  data() {
    return {
      loading: false,
      user: { events: [], playlists: [], tracks: [] },
      errorMessage: "",
      selectedTrackURI: "spotify:track:3lzUeaCbcCDB5IXYfqWRlF", // Updated to null initially
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
        const response = await axios.get(`${API_BASE_URL}/tracks/top`, {
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
    setSelectedTrackURI(trackURI) {
      this.selectedTrackURI = trackURI; // Update the selected track URI
    },
    incrementPlayCount(trackId) {
      incrementTrackPlayCount(trackId)
        .then(() => {
          console.log(`Play count incremented for track ID: ${trackId}`);
        })
        .catch((error) => {
          console.error(`Error incrementing play count for track ID: ${trackId}`, error);
        });
    },
  },
  mounted() {
    this.fetchEvents();
  },
};
</script>
