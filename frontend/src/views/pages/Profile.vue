<!-- Profile.vue -->
<template>
  <div class="profile-page" style="max-width: 1000px; margin: 2rem auto">
    <div
      v-if="loading"
      class="p-d-flex p-jc-center p-ai-center"
      style="font-size: 1.2rem; padding: 2rem"
    >
      <span>Loading…</span>
    </div>

    <div v-else-if="errorMessage" class="p-error p-4 text-red-500">
      {{ errorMessage }}
    </div>

    <div v-else class="profile-details p-card p-p-4 p-shadow-4 pt-4">
      <img
        :src="user.profile_photo_url || '/public/profile.svg'"
        alt="Profile Photo"
        style="
          width: 120px;
          height: 120px;
          object-fit: cover;
          border-radius: 50%;
          border: 3px solid var(--primary-color);
          display: block;
          margin: 0 auto;
        "
      />
      <h1 class="p-mt-3 text-center">{{ user.user_name }}</h1>
      <!-- <p><strong>Email:</strong> {{ user.email_address }}</p>
      <p><strong>Phone:</strong> {{ user.phone_number }}</p>
      <p><strong>Location:</strong> {{ user.location }}</p> -->

      <template v-if="hasContent">
        <section class="p-4" v-if="user.events.length">
          <h2 class="text-xl font-semibold mb-3">Events</h2>
          <div class="flex space-x-4 overflow-x-auto pb-4">
            <EventCard v-for="e in user.events" :key="e.event_id" :event="e" />
          </div>
        </section>

        <section class="p-4" v-if="user.playlists.length">
          <h2 class="text-xl font-semibold mb-3">Playlists</h2>
          <div class="flex space-x-4 overflow-x-auto pb-4">
            <PlaylistCard
              v-for="p in user.playlists"
              :key="p.playlist_id"
              :playlist="p"
            />
          </div>
        </section>

        <section class="p-4" v-if="user.tracks.length">
          <h2 class="text-xl font-semibold mb-3">Tracks</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            <TrackCard v-for="t in user.tracks" :key="t.track_id" :track="t" />
          </div>
        </section>
      </template>
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
            EventService.getEventsByUserId(userId),
            PlaylistService.getPlaylistsByUserId(userId),
            getUserTracksById(userId),
          ]);
          this.user = { ...u, events, playlists, tracks };
        } else {
          // current logged‑in user
          const [u, events, playlists, tracks] = await Promise.all([
            UserService.getUser({ withCredentials: true }), // /me endpoint inside UserService
            EventService.getEventsForCurrentUser(),
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
