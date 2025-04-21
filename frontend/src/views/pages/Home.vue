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
          <div class="p-6 bg-gradient-to-r from-green-400 via-green-500 to-green-600 rounded-lg shadow-lg"
            v-if="user.tracks.length">
            <!-- Title with Emoji -->
            <h2 class="text-3xl font-extrabold text-center text-gray-800 mb-4 animate-heading">
            Top Played Tracks
          </h2>

          <!-- Subheading with Emoji -->
          <p class="text-center text-lg text-gray-700 mb-6 opacity-90 animate-subheading">
            Your most played tracks, right here!
          </p>

            <!-- Track Cards Grid -->
            <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 gap-4 auto-rows-fr">
              <!-- Track Card Loop -->
              <TrackCard v-for="track in user.tracks" :key="track.track_id" :track="track" :state="'redirect'"
                @track-selected="setSelectedTrackURI" @click="handleTrackFunction(track.track_id)"
                :selectedTracks="selectedTracks"  :liked="likedTracks.includes(track.track_id)"  @track-unliked="handleTrackUnliked"
                @track-liked="handleTrackLiked" 
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
        <div class v-if="filter === 'all' || filter === 'music'">
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
    <SpotifyPlayer />
  </div>
</template>

<script>
import axios from "axios";
import { API_BASE_URL } from "@/service/apiConfig";
import EventService from "@/service/EventService";
import TrackService from "@/service/TrackService";
import EventCard from "@/components/EventCard.vue";
import TrackCard from "@/components/TrackCard.vue";
import PlaylistCard from "@/components/PlaylistCard.vue";
import SpotifyPlayer from "@/components/SpotifyPlayer.vue";
import RecommendedTracks from "@/components/RecommendedTracks.vue";
import { incrementTrackPlayCount } from "@/service/TrackService";
import { useSpotifyStore } from "@/store/SpotifyStore";

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
      likedTracks: [],
      selectedTracks: [],
      errorMessage: "",
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
  async handleTrackFunction(trackId) {
    try {
      await incrementTrackPlayCount(trackId);
      await this.refreshLikedTracks();

      const isLiked = this.likedTracks.includes(trackId);

      console.log("likedTracks", this.likedTracks);

      if (isLiked) {
        await this.handleTrackUnliked(trackId);
      } else {
        await this.handleTrackLiked(trackId);
      }

      // Optional double-check to ensure liked tracks are updated
      await this.refreshLikedTracks();

      // Emit event with track ID
      this.$emit("track-selected", trackId);
    } catch (error) {
      console.error("Error handling track:", error);
      this.$toast.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Something went wrong while updating the track. Please try again.',
        life: 3000,
      });
    }
  },

  async refreshLikedTracks() {
  try {
    const response = await TrackService.likedTracks();
    if (Array.isArray(response)) {
      this.likedTracks = response.map(t => t.track_id);
      console.log("Updated likedTracks:", this.likedTracks); // Add this log
    } else {
      console.error("Invalid liked tracks data format");
      this.$toast.add({
        severity: 'warn',
        summary: 'Warning',
        detail: 'Could not fetch liked tracks properly. Please try again.',
        life: 3000,
      });
    }
  } catch (error) {
    console.error("Error fetching liked tracks:", error);
    this.$toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Unable to fetch liked tracks. Please check your connection.',
      life: 3000,
    });
  }
},


  async handleTrackLiked(trackId) {
    try {
      await TrackService.likeTrack(trackId);
      this.$toast.add({
        severity: 'success',
        summary: 'Track Liked',
        detail: 'This track has been added to your liked tracks!',
        life: 3000,
      });
    } catch (error) {
      console.error("Error liking track:", error);
      this.$toast.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Failed to like the track. Please try again.',
        life: 3000,
      });
    }
  },

  async handleTrackUnliked(trackId) {
    try {
      await TrackService.unlikeTrack(trackId);
      this.$toast.add({
        severity: 'info',
        summary: 'Track Unliked',
        detail: 'This track has been removed from your liked tracks.',
        life: 3000,
      });
    } catch (error) {
      console.error("Error unliking track:", error);
      this.$toast.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Failed to unlike the track. Please try again.',
        life: 3000,
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
  },

  async mounted() {
  this.loading = true;
  await Promise.all([
    this.getTracksForUser(),
    this.getPlaylistsForUser(),
    this.getEventsByUserId(),
    this.getEvents(),
    this.refreshLikedTracks(),
  ]);
  this.loading = false;
},
};


</script>

<style scoped>

@keyframes bounceIn {
  0% {
    transform: scale(0.8);
    opacity: 0;
  }
  50% {
    transform: scale(1.1);
    opacity: 1;

  }
  100% {
    transform: scale(1);
    opacity: 1;

  }
}

@keyframes colorPop {
  0% {
    color: #000000; /* Initial color */
  }
  50% {
    color: #000000; /* Midway color */
    text-shadow: 0 0 2px #ffeb3b, 0 0 10px #ffeb3b;
  
  }
  100% {
    color: #000000; /* Final color */
  }
}

.animate-heading {
  animation: bounceIn 1.2s ease-out, colorPop 2s ease-in-out ;
}

.animate-subheading {
  animation: bounceIn 1.5s ease-out, colorPop 2s ease-in-out ;
}


</style>