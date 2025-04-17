<template>
  <div
    class="profile-page"
    style="max-width: 1000px; margin: 2rem auto; text-align: center"
  >
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




    <div v-else class="profile-details p-card p-p-4 p-shadow-4 m-4 p-8">

      <img
        :src="user.profile_photo_url || '/public/profile.svg'"
        alt="Profile Photo"
        style="
          width: 120px;
          height: 120px;
          object-fit: cover;
          border-radius: 50%;
          border: 3px solid var(--primary-color);
          display: block; /* makes margin work */
          margin: 0 auto; /* centres horizontally */
        "
      />
      <br> 
      
      <div class="p-d-flex p-flex-column">
  <h2 class="text-xl font-bold">{{ user.user_name }}</h2>
  <p class="text-sm text-muted p-mt-1">
    {{ user.email_address }} · {{ user.phone_number }} · {{ user.location }}
  </p>
</div>
</div>



    <br>

      <template v-if="hasContent">
        <section class="p-4" v-if="user.events.length">
          <h2 class="text-xl font-semibold mb-3 text-left">Liked Events</h2>

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
        </section>

        <!-- <section class="p-4" v-if="user.playlists.length">
          <h2 class="text-xl font-semibold mb-3">Playlists</h2>
          <div class="flex space-x-4 overflow-x-auto pb-4">
            <PlaylistCard
              v-for="p in user.playlists"
              :key="p.playlist_id"
              :playlist="p"
            />
          </div>
        </section> -->

        <section class="p-4" v-if="user.tracks.length">
          <h2 class="text-xl font-semibold mb-3 text-left">Liked Tracks</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            <TrackCard v-for="t in user.tracks" :key="t.track_id" :track="t" />
          </div>
        </section>

        <!-- <section class="p-4" v-if="user.tracks.length">
          <h2 class="text-xl font-semibold mb-3">Artists</h2>
          <SpotifyPlayer />
        </section> -->
      </template>
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

      // First, get the core user data. If this fails, the page cannot render.
      let userData;
      try {
        userData = await UserService.getUser({ withCredentials: true });
      } catch (err) {
        console.error("Failed to fetch user details:", err);
        this.errorMessage =
          "Failed to fetch user details. Please try again later.";
        this.loading = false;
        return;
      }

      // Now load additional data; if any fail, log the error and fallback to empty lists.
      const eventsPromise = userId
        ? EventService.getEventsByUserId(userId)
        : EventService.getEventsForCurrentUser();
      const playlistsPromise = userId
        ? PlaylistService.getPlaylistsByUserId(userId)
        : PlaylistService.getPlaylistsForUser();
      const tracksPromise = userId
        ? getUserTracksById(userId)
        : getUserTracks();

      const results = await Promise.allSettled([
        eventsPromise,
        playlistsPromise,
        tracksPromise,
      ]);
      const events = results[0].status === "fulfilled" ? results[0].value : [];
      const playlists =
        results[1].status === "fulfilled" ? results[1].value : [];
      const tracks = results[2].status === "fulfilled" ? results[2].value : [];

      if (results.some((result) => result.status === "rejected")) {
        console.error("One or more profile sections failed to load", results);
        // You might choose to provide more detail, but a generic message keeps things simple.
        this.errorMessage = "Some parts of your profile failed to load.";
      } else {
        this.errorMessage = "";
      }

      this.user = { ...userData, events, playlists, tracks };
      this.loading = false;
    },
  },
  mounted() {
    this.fetchData();
  },
};
</script>
