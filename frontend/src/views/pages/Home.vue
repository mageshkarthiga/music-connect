<template>
  <div>
    <!-- Top Bar with Filter Buttons -->
    <div class="p-4 flex space-x-4">
      <Button label="All" severity="secondary" outlined @click="filterContent('all')" />
      <Button label="Music" severity="secondary" outlined @click="filterContent('music')" />
      <Button label="Events" severity="secondary" outlined @click="filterContent('events')" />
    </div>

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
        <!-- Tracks -->
        <div class="p-4" v-if="filter === 'all' || filter === 'music'">
          <h2 class="text-xl font-semibold mb-3"> Frequently Accessed Tracks</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            <!-- Handle track-selected event -->

                      <TrackCard
            v-for="track in user.tracks"
            :key="track.track_id"
            :track="track"
            :state="'redirect'"
            :liked="likedTracks.includes(track.track_id)"
            :selectedTracks="selectedTracks"
            @track-selected="setSelectedTrackURI"
          />

          <!-- <TrackCard
            v-for="track in user.tracks"
            :key="track.track_id"
            :track="track"
            :state="'redirect'"
            :liked="likedTracks.includes(track.track_id)"
            :selectedTracks="selectedTracks"
            @track-liked="handleTrackLiked"
            @track-unliked="handleTrackUnliked"
            @track-selected="setSelectedTrackURI"
          /> -->

          </div>
        </div>

        <br>

        <!-- User's liked events -->
        <!-- <div v-if="(filter === 'all' || filter === 'events') && user.events.length">
          <h2 class="text-xl font-semibold mb-3">Liked Events</h2>
          <div class="flex space-x-4 overflow-x-auto pb-4">
            <EventCard
              v-for="event in user.events"
              :key="event.event_id"
              :event="event"
              :liked="true"
              @event-unliked="handleEventUnliked"
              @event-liked="handleEventLiked"
            />
          </div>
        </div> -->

        <br>

        <!-- Discoverable events -->
        <div v-if="(filter === 'all' || filter === 'events') && otherEvents.length">
          <h2 class="text-xl font-semibold mb-3">Discover Events</h2>
          <div class="flex space-x-4 overflow-x-auto pb-4">
            <EventCard
              v-for="event in otherEvents"
              :key="event.event_id"
              :event="event"
              :liked="false"
              @event-unliked="handleEventUnliked"
              @event-liked="handleEventLiked"
            />
          </div>
        </div>

        <div class="p-4" v-if="filter === 'all' || filter === 'music'">
          <div class="font-semibold text-xl mb-4">Recommended music </div>
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
import EventService from "@/service/EventService";
import EventCard from "@/components/EventCard.vue";
import TrackCard from "@/components/TrackCard.vue";
import PlaylistCard from "@/components/PlaylistCard.vue";
import SpotifyPlayer from "@/components/SpotifyPlayer.vue";
import RecommendedTracks from "@/components/RecommendedTracks.vue";

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
      events: [],
      likedTracks: [],
      selectedTracks: [],
      errorMessage: "",
      selectedTrackURI: "spotify:track:3lzUeaCbcCDB5IXYfqWRlF", // Updated to null initially
      filter: 'all',  // Default filter value
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

    async handleTrackLiked(trackId) {
      const likedTrack = this.user.tracks.find(t => t.track_id === trackId);
      if (likedTrack) {
        this.likedTracks.push(likedTrack.track_id);
        this.user.tracks = this.user.tracks.filter(t => t.track_id !== trackId);
      }

      this.$toast.add({
        severity: 'success',
        summary: 'Track Liked',
        detail: 'This track has been added to your liked tracks!',
        life: 3000,
      });
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
        life: 3000,
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
        const response = await axios.get(`${API_BASE_URL}/events`, {
          withCredentials: true,
        });
        this.events = response.data;
      } catch (err) {
        this.handleError(err, "events");
      }
    },

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
          this.getEvents(),
          this.getEventsByUserId(),
          this.getPlaylistsForUser(),
          this.getTracksForUser(),
        ]);
      } finally {
        this.loading = false;
      }
    },

    setSelectedTrackURI(trackURI) {
      this.selectedTrackURI = trackURI;
    },

    // Filter Content based on selected category
    filterContent(category) {
      this.filter = category;
    }
  },

  mounted() {
    this.fetchEvents();
  },
};
</script>
