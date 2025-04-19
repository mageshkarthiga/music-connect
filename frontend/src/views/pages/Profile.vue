<template>
  <div
    class="max-w-screen-md mx-auto my-8 bg-surface-0 dark:bg-surface-900 rounded-lg shadow-lg text-surface-900 dark:text-white">
    <!-- Loading Spinner -->
    <div v-if="loading" class="flex justify-center items-center py-10">
      <ProgressSpinner style="width: 50px; height: 50px;" strokeWidth="5" animationDuration=".7s" />
    </div>
    <!-- Error Message -->
    <div v-else-if="errorMessage" class="p-error p-4 text-red-500">
      {{ errorMessage }}
    </div>

    <!-- Profile Details -->
    <div v-else class="profile-details p-card p-p-4 p-shadow-4 mt-4 p-8">
      <img :src="user?.profile_photo_url || '/profile.svg'" alt="Profile Photo"
        class="w-[120px] h-[120px] object-cover rounded-full border-4 border-primary" />
      <br>
      <div class="p-d-flex p-flex-column">
        <h2 class="text-xl font-bold">{{ user.user_name }}</h2>
        <p class="text-sm text-muted p-mt-1">
          {{ user.email_address }} · {{ user.phone_number }} · {{ user.location }}
        </p>
      </div>
    </div>

    <br>
    <!-- Pending Friend Requests -->
    <div class="p-4">
      <h2 class="text-xl font-semibold mb-3 text-left">Pending Friend Requests</h2>
      <div class="flex space-x-4 overflow-x-auto pb-4">
        <div v-for="u in user.pendingFriendRequests" :key="u.user_id" class="min-w-[280px] max-w-md">
          <UserCard :user="u" :isPending="true" @accept="handleAccept" @reject="handleReject" />
        </div>
      </div>
    </div>
    <br>

    <!-- Content Sections -->
    <template v-if="hasContent">
      <!-- Liked Events -->
      <section v-if="user.events.length" class="p-4">
        <h2 class="text-xl font-semibold mb-3 text-left">Liked Events</h2>
        <div class="flex space-x-4 overflow-x-auto pb-4">
          <EventCard v-for="event in user.events" :key="event.event_id" :event="event" :liked="true"
            @event-unliked="handleEventUnliked" @event-liked="handleEventLiked" />
        </div>
      </section>

      <!-- Tracks -->
      <section v-if="user.tracks.length" class="p-4">
        <h2 class="text-xl font-semibold mb-3 p-5">Tracks</h2>
        <div class="flex space-x-4 overflow-x-auto pb-4">
          <div v-for="t in user.tracks" :key="t.track_id" class="min-w-[280px] max-w-md">
            <TrackCard :track="t" />
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<script>
import EventCard from "@/components/EventCard.vue";
import PlaylistCard from "@/components/PlaylistCard.vue";
import TrackCard from "@/components/TrackCard.vue";
import SpotifyPlayer from "@/components/SpotifyPlayer.vue";
import UserCard from "@/components/UserCard.vue";
import { getPendingFriendRequests } from "@/service/FriendService"
import UserService from "@/service/UserService";
import EventService from "@/service/EventService";
import PlaylistService from "@/service/PlaylistService";
import {
  getUserTracksById,
  getUserTracks,
  getFavUserTracksById,
  getFavUserTracks,
} from "@/service/TrackService";

export default {
  name: "Profile",
  components: { EventCard, PlaylistCard, TrackCard, SpotifyPlayer },

  data: () => ({
    loading: true,
    user: { events: [], playlists: [], tracks: [], friendRequests:[] },
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
          const [u, events, playlists, tracks, friendRequests] = await Promise.all([
            UserService.getUserByUserId(userId),
            EventService.getFavEventsByUserId(userId),
            PlaylistService.getPlaylistsByUserId(userId),
            getFavUserTracksById(userId),
            getPendingFriendRequests(userId),
          ]);
          this.user = { ...u, events, playlists, tracks, friendRequests };
        } else {
          // current logged‑in user
          const [u, events, playlists, tracks] = await Promise.all([
            UserService.getUser({ withCredentials: true }), // /me endpoint inside UserService
            EventService.getFavEventsForCurrentUser(),
            PlaylistService.getPlaylistsForUser(),
            getFavUserTracks(),
            getPendingFriendRequests(),
          ]);
          this.user = { ...u, events, playlists, tracks, friendRequests };
          console.log("user", this.user);
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
