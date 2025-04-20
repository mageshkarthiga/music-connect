<template>
  <div
    class="max-w-screen-md mx-auto my-8 bg-surface-0 dark:bg-surface-900 rounded-lg shadow-lg text-surface-900 dark:text-white p-3">
    <!-- Loading Spinner -->
    <div v-if="loading" class="flex justify-center items-center py-10">
      <ProgressSpinner style="width: 50px; height: 50px;" strokeWidth="5" animationDuration=".7s" />
    </div>
    <!-- Error Message -->
    <div v-else-if="errorMessage" class="p-error p-4 text-red-500">
      {{ errorMessage }}
    </div>

    <!-- Profile Details -->
    <div v-else class="profile-details p-shadow-4 mt-4 p-8">
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
    <Divider />
    <br>

    <!-- Content Sections -->
    <template v-if="hasContent">
      <!-- Pending Friend Requests -->
      <section v-if="user.friendRequests?.users?.length" class="p-4">
        <h2 class="text-xl font-semibold mb-3 text-left">Pending Friend Requests</h2>
        <div class="flex space-x-4 overflow-x-auto pb-4">
          <div v-for="u in user.friendRequests.users" :key="u.user_id" class="min-w-[280px] max-w-md">
            <UserCard :user="u" @accept="handleAccept(u.user_id)" @reject="handleReject(u.user_id)" :accept="true" :reject="true"/>
          </div>
        </div>
      </section>
      <section v-if="user.friends.length" class="p-4">
        <h2 class="text-xl font-semibold mb-3 text-left">Friends</h2>
        <div class="flex space-x-4 overflow-x-auto pb-4">
          <div v-for="u in user.friends" :key="u.user_id" class="min-w-[280px] max-w-md">
            <UserCard :user="u" :remove="true" @remove="handleRemove(u.user_id)"/>
          </div>
        </div>
        <Divider />
      </section>
      <!-- Liked Events -->
      <section v-if="user.events.length" class="p-4">
        <h2 class="text-xl font-semibold mb-3 text-left">Liked Events</h2>
        <div class="flex space-x-4 overflow-x-auto pb-4">
          <EventCard v-for="event in user.events" :key="event.event_id" :event="event" :liked="true"
            @event-unliked="handleEventUnliked" @event-liked="handleEventLiked" />
        </div>
      </section>

      <!-- Tracks -->
        <section class="p-4" v-if="user.tracks.length">
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
import UserCard from "@/components/UserCard.vue";
import UserService from "@/service/UserService";
import EventService from "@/service/EventService";
import PlaylistService from "@/service/PlaylistService";
import {
  getUserTracksById,
  getUserTracks,
  getFavUserTracksById,
  getFavUserTracks,
} from "@/service/TrackService";
import FriendService from "@/service/FriendService";

export default {
  name: "Profile",
  components: { EventCard, PlaylistCard, TrackCard, UserCard },

  data: () => ({
    loading: true,
    user: { events: [], playlists: [], tracks: [], friends: [], friendRequests: [] },
    errorMessage: "",
  }),

  computed: {
    hasContent() {
      const { events, playlists, tracks } = this.user;
      return events.length || playlists.length || tracks.length ;
    },
  },

  methods: {
    async fetchData() {
      this.loading = true;
      const userId = Number(this.$route.query.user_id);

      try {
        if (!Number.isNaN(userId)) {
          // Fetch data for an explicit user (not the logged-in user)
          const [u, events, playlists, tracks, friends] = await Promise.all([
            UserService.getUserByUserId(userId),
            EventService.getFavEventsByUserId(userId),
            PlaylistService.getPlaylistsByUserId(userId),
            getFavUserTracksById(userId),
            FriendService.getFriends()
          ]);
          this.user = { ...u, events, playlists, tracks, friendRequests: null }; // No friendRequests for other users
        } else {
          // Fetch data for the logged-in user
          const [u, events, playlists, tracks, friends, friendRequests] = await Promise.all([
            UserService.getUser({ withCredentials: true }),
            EventService.getFavEventsForCurrentUser(),
            PlaylistService.getPlaylistsForUser(),
            getFavUserTracks(),
            FriendService.getFriends(),
            FriendService.getPendingFriendRequests(),
          ]);
          this.user = { ...u, events, playlists, tracks, friends, friendRequests };
        }
      } catch (err) {
        console.error("profile fetch error:", err);
        this.errorMessage =
          err?.response?.data?.message ?? "Failed to fetch profile.";
      } finally {
        this.loading = false;
      }
    },
    handleAccept(userId) {
      FriendService.acceptFriendRequest(userId)
        .then(() => {
          this.user.friendRequests.users = this.user.friendRequests.users.filter(
            (u) => u.user_id !== userId
          );
          this.user.friends.push(this.user.friendRequests.users.find((u) => u.user_id === userId));
          this.$toast.add({
            severity: "success",
            summary: "Success",
            detail: "Friend request accepted.",
          });
        })
        .catch((err) => {
          console.error("Error accepting friend request:", err);
          this.$toast.add({
            severity: "error",
            summary: "Error",
            detail: "Failed to accept friend request.",
          });
        });
    },
    handleReject(userId) {
      FriendService.rejectFriendRequest(userId)
        .then(() => {
          this.user.friendRequests.users = this.user.friendRequests.users.filter(
            (u) => u.user_id !== userId
          );
          this.$toast.add({
            severity: "info",
            summary: "Success",
            detail: "Friend request rejected.",
          });
        })
        .catch((err) => {
          console.error("Error rejecting friend request:", err);
          this.$toast.add({
            severity: "error",
            summary: "Error",
            detail: "Failed to reject friend request.",
          });
        });
    },
    handleRemove(userId) {
      FriendService.removeFriend(userId)
        .then(() => {
          this.user.friends = this.user.friends.filter((u) => u.user_id !== userId);
          this.$toast.add({
            severity: "info",
            summary: "Success",
            detail: "Friend removed.",
          });
        })
        .catch((err) => {
          console.error("Error removing friend:", err);
          this.$toast.add({
            severity: "error",
            summary: "Error",
            detail: "Failed to remove friend.",
          });
        });
    },
  },

  mounted() {
    this.fetchData();
  },
};
</script>
