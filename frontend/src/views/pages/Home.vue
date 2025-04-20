<template>
  <div>
    <!-- Top Bar with Filter Buttons -->
    <div class="p-4 flex flex-wrap gap-3 justify-center sm:justify-start">
      <Button label="All" icon="pi pi-list" :severity="filter === 'all' ? 'success' : 'secondary'"
        :outlined="filter !== 'all'" class="transition-all duration-200 rounded-full px-4 py-2"
        @click="filterContent('all')" />
      <Button label="Music" icon="pi pi-headphones" :severity="filter === 'music' ? 'success' : 'secondary'"
        :outlined="filter !== 'music'" class="transition-all duration-200 rounded-full px-4 py-2"
        @click="filterContent('music')" />
      <Button label="Events" icon="pi pi-calendar" :severity="filter === 'events' ? 'success' : 'secondary'"
        :outlined="filter !== 'events'" class="transition-all duration-200 rounded-full px-4 py-2"
        @click="filterContent('events')" />
    </div>

    <!-- Loading Spinner -->
    <div v-if="loading" class="flex justify-center items-center py-10">
      <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="5" animationDuration=".7s" />
    </div>

    <!-- Error Message -->
    <div v-if="errorMessage" class="p-4 text-red-500">
      {{ errorMessage }}
    </div>

    <!-- Content -->
    <div v-if="!loading">
      <template v-if="hasContent">
        <!-- Tracks -->
        <div class="p-4" v-if="filter === 'all' || filter === 'music'">
          <!-- Tracks Section -->
          <div class="p-6 bg-gradient-to-r from-green-400 via-green-400 to-green-500 rounded-lg shadow-lg"
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
            <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 gap-4 auto-rows-fr">
              <!-- Track Card Loop -->
              <TrackCard v-for="track in user.tracks" :key="track.track_id" :track="track" :state="'redirect'"
                @track-selected="setSelectedTrackURI"
                :selectedTracks="selectedTracks"
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
            class="p-6 bg-gradient-to-r from-green-400 via-green-400 to-green-500 rounded-lg shadow-xl text-center">
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
        </div>
        <br>

        <br>

        <!-- Discoverable events -->
        <div v-if="(filter === 'all' || filter === 'events') && otherEvents.length" class="mb-10 p-3">
          <h2 class="text-xl font-semibold mb-3">Discover Events</h2>
          <div class="flex space-x-4 overflow-x-auto pb-4">
            <EventCard v-for="event in otherEvents" :key="event.event_id" :event="event" :liked="false"
              @event-unliked="handleEventUnliked" @event-liked="handleEventLiked" />
          </div>
        </div>

        <br />

        <!-- Recommended Music -->
        <div class="p-4" v-if="filter === 'all' || filter === 'music'">
          <div class="font-semibold text-xl mb-4">Recommended music </div>
          <RecommendedTracks @track-selected="setSelectedTrackURI" />
        </div>
      </template>

      <!-- No Content -->
      <template v-else>
        <div class="p-4">
          <p>No events, playlists, or tracks found.</p>
        </div>
      </template>
    </div>

    <!-- Spotify Player -->
    <SpotifyPlayer />
  </div>
</template>

<script>
import axios from "axios";
import EventService from "@/service/EventService";
import EventCard from "@/components/EventCard.vue";
import TrackCard from "@/components/TrackCard.vue";
import SpotifyPlayer from "@/components/SpotifyPlayer.vue";
import RecommendedTracks from "@/components/RecommendedTracks.vue";
import { incrementTrackPlayCount } from "@/service/TrackService";
import { useSpotifyStore } from "@/store/SpotifyStore";
import { API_BASE_URL } from "@/service/apiConfig";

export default {
  components: {
    EventCard,
    TrackCard,
    SpotifyPlayer,
    RecommendedTracks,
  },

  
  data() {
    return {
      loading: false,
      user: { events: [], playlists: [], tracks: [] },
      events: [],
      likedTrackIds: [],

      selectedTracks: [],
      errorMessage: "",
      filter: 'all',  // Default filter value
      API_BASE_URL,
    };
  },

  
  computed: {
    hasContent() {
      return (
        this.user.events.length ||
        this.events.length ||
        this.user.playlists.length ||
        this.user.tracks.length
      );
    },

    
    otherEvents() {
      const userEventIds = new Set(this.user.events.map(e => e.event_id));
      return this.events.filter(e => !userEventIds.has(e.event_id));
    }
  },
  methods: {

    handleTrackClick(trackId) {
      this.incrementPlayCount(trackId);
      if (this.isLiked(trackId)) {
        this.handleTrackUnliked(trackId);
      } else {
        this.handleTrackLiked(trackId);
      }
    },

    async handleTrackLiked(trackId) {
      const likedTrack = this.user.tracks.find(t => t.track_id === trackId);
      if (likedTrack) {
        this.likedTrackIds.push(likedTrack.track_id);

        this.user.tracks = this.user.tracks.filter(t => t.track_id !== trackId);
      }

      this.$toast.add({
        severity: 'success',
        summary: 'Track Liked',
        detail: 'This track has been added to your liked tracks!',
        life: 3000,
      });
    },

    async handleTrackUnliked(trackId) {
      const unlikedTrack = this.user.tracks.find(t => t.track_id === trackId);

      if (unlikedTrack) {
        // Call the API to unlike the track on the backend
        await EventService.unlikeTrack(trackId);

        // Remove the track from user.tracks (liked tracks)
        this.user.tracks = this.user.tracks.filter(t => t.track_id !== trackId);

      this.$toast.add({
        severity: 'info',
        summary: 'Track Unliked',
        detail: 'This track has been removed from your liked tracks.',

      });
      }
    },

    async handleEventLiked(eventId) {
      const likedEvent = this.otherEvents.find(e => e.event_id === eventId);
      if (likedEvent) {
        this.user.events.push(likedEvent);
        this.events = this.events.filter(e => e.event_id !== eventId);
      }

      this.$toast.add({
        severity: 'success',
        summary: 'Event Liked',
        detail: 'This event has been added to your liked events!',
     
      });
    },

    async handleEventUnliked(eventId) {
      const unlikedEvent = this.user.events.find(e => e.event_id === eventId);

      if (unlikedEvent) {
        // Call the API to unlike the event on the backend
        await EventService.unlikeEvent(eventId);

        // Remove the event from user.events (liked events)
        this.user.events = this.user.events.filter(e => e.event_id !== eventId);

        // Add it back to discoverable events if not already there
        if (!this.events.some(e => e.event_id === eventId)) {
          this.events.push(unlikedEvent);
        }
      }

      this.$toast.add({
        severity: 'warn',
        summary: 'Event Unliked',
        detail: 'This event has been removed from your liked events.',
        life: 3000,
      });
    },

    async getEvents() {
      try {
        const response = await axios.get(`${this.API_BASE_URL}/events`, {
          withCredentials: true,
        });
        this.events = response.data;
      } catch (err) {
        this.handleError(err, "events");
      }
    },

    async getEventsByUserId() {
      try {
        const response = await axios.get(`${this.API_BASE_URL}/me/events`, {
          withCredentials: true,
        });
        this.user.events = Array.isArray(response.data) ? response.data : [];
      } catch (err) {
        this.handleError(err, "events");
      }
    },

    async getPlaylistsForUser() {
      try {
        const response = await axios.get(`${this.API_BASE_URL}/me/playlists`, {
          withCredentials: true,
        });
        this.user.playlists = response.data;
      } catch (err) {
        this.handleError(err, "playlists");
      }
    },

    async getTracksForUser() {
      try {
        const response = await axios.get(`${this.API_BASE_URL}/tracks/top`, {
          withCredentials: true,
        });
        this.user.tracks = response.data;
      } catch (err) {
        this.handleError(err, "tracks");
      }
    },

    async getLikedTracks() {
      try {
        const response = await axios.get(`${API_BASE_URL}/likedTracks`, {
          withCredentials: true,
        });

        // this.likedTracks = response.data;
        // this.likedTrackIds = response.data.likedTracks.map(track => track.track_id);

        const likedTracks = Array.isArray(response.data)
          ? response.data
          : response.data?.likedTracks ?? [];

        this.likedTracks = likedTracks;
        this.likedTrackIds = likedTracks.map(track => track.track_id);
      } catch (err) {
        this.handleError(err, "liked tracks");
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
          this.getEvents(),
          this.getEventsByUserId(),
          this.getPlaylistsForUser(),
          this.getTracksForUser(),
          this.getLikedTracks(),
        ]);
      } finally {
        this.loading = false;
      }
    },
    setSelectedTrackURI(trackURI) {
      const spotifyStore = useSpotifyStore();
      spotifyStore.spotifyUri = trackURI;
      // this.selectedTrackURI = trackURI;
    },
    // Filter Content based on selected category
    filterContent(category) {
      this.filter = category;
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
    isLiked(trackId) {
      return this.likedTrackIds.includes(trackId);
    },
  },

  mounted() {
    this.fetchEvents();
    this.getLikedTracks();
  },
};
</script>
