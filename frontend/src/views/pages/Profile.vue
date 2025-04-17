<template>
  <div
    class="max-w-screen-md mx-auto my-8 bg-surface-0 dark:bg-surface-900 rounded-lg shadow-lg text-surface-900 dark:text-white"
  >
    <div v-if="loading" class="flex justify-center items-center text-lg p-8">
      <span>Loading…</span>
    </div>

    <div
      v-else
      class="bg-surface-card border border-surface-border rounded-lg shadow-md text-center flex flex-col items-center px-6 py-10 min-h-[300px] text-surface-900 dark:text-white"
    >
      <img
        :src="user?.profile_photo_url || '/public/profile.svg'"
        alt="Profile Photo"
        class="w-[120px] h-[120px] object-cover rounded-full border-4 border-primary"
      />
      <h1 class="mt-4 text-xl font-semibold">
        {{ user?.user_name || "Unknown User" }}
      </h1>

      <div
        v-if="errorMessage && !hasContent"
        class="text-red-600 dark:text-red-400 mt-6"
      >
        <p>{{ errorMessage }}</p>
        <p class="mt-2 text-surface-600 dark:text-surface-300">
          No data found for this user.
        </p>
      </div>

      <!-- Events Section -->
      <section class="w-full p-4 text-left">
        <h2 class="text-2xl font-semibold mb-3">Favourite Events</h2>
        <div
          v-if="user.events.length > 0 && user.events[0].event_id"
          class="flex space-x-4 overflow-x-auto pb-4"
        >
          <EventCard v-for="e in user.events" :key="e.event_id" :event="e" />
        </div>
        <p v-else class="text-surface-600 dark:text-surface-300">
          No events to show.
        </p>
      </section>

      <!-- Playlists Section -->
      <section class="w-full p-4 text-left">
        <h2 class="text-2xl font-semibold mb-3">Playlists</h2>
        <div
          v-if="user.playlists.length"
          class="flex space-x-4 overflow-x-auto pb-4"
        >
          <PlaylistCard
            v-for="p in user.playlists"
            :key="p.playlist_id"
            :playlist="p"
          />
        </div>
        <p v-else class="text-surface-600 dark:text-surface-300">
          No playlists found.
        </p>
      </section>

      <!-- Tracks Section -->
      <section class="w-full p-4 text-left">
        <h2 class="text-xl font-semibold mb-3">Favourite Tracks</h2>

        <div
          v-if="user.tracks.length"
          class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4"
        >
          <TrackCard
            v-for="track in user.tracks"
            :key="track.track_id"
            :track="track"
            :state="'redirect'"
            @track-selected="setSelectedTrackURI"
          />
        </div>

        <p v-else class="text-surface-600 dark:text-surface-300">
          No tracks available.
        </p>
      </section>
    </div>
  </div>
</template>

<script>
import EventCard from "@/components/EventCard.vue";
import PlaylistCard from "@/components/PlaylistCard.vue";
import TrackCard from "@/components/TrackCard.vue";
import SpotifyPlayer from "@/components/SpotifyPlayer.vue";

import UserService from "@/service/UserService";
import EventService from "@/service/EventService";
import PlaylistService from "@/service/PlaylistService";
import { getUserTracksById, getUserTracks } from "@/service/TrackService";

export default {
  name: "Profile",
  components: { EventCard, PlaylistCard, TrackCard, SpotifyPlayer },

  data: () => ({
    loading: true,
    user: { events: [], playlists: [], tracks: [] },
    errorMessage: "",
  }),

  computed: {
    hasContent() {
      const { events, playlists, tracks } = this.user;
      return events.length || playlists.length || tracks.length;
    },
  },

  methods: {
    async fetchData() {
      this.loading = true;
      const userId = Number(this.$route.query.user_id);

      try {
        if (!Number.isNaN(userId)) {
          // explicit user
          const [u, events, playlists, tracks] = await Promise.all([
            UserService.getUserByUserId(userId),
            EventService.getFavEventsByUserId(userId),
            PlaylistService.getPlaylistsByUserId(userId),
            getUserTracksById(userId),
          ]);
          this.user = { ...u, events, playlists, tracks };
        } else {
          // current logged‑in user
          const [u, events, playlists, tracks] = await Promise.all([
            UserService.getUser({ withCredentials: true }), // /me endpoint inside UserService
            EventService.getFavEventsForCurrentUser(),
            PlaylistService.getPlaylistsForUser(),
            getUserTracks(),
          ]);
          this.user = { ...u, events, playlists, tracks };
        }
      } catch (err) {
        console.error("profile fetch error:", err);
        this.errorMessage =
          err?.response?.data?.message ?? "Failed to fetch profile.";
      } finally {
        this.loading = false;
      }
    },
  },

  mounted() {
    this.fetchData();
  },
};
</script>
