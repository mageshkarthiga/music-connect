<template>
  <div
    class="max-w-screen-xl mx-auto my-8 bg-surface-0 dark:bg-surface-900 rounded-lg shadow-lg text-surface-900 dark:text-white p-4">

    <!-- Loading Spinner -->
    <div v-if="loading" class="flex justify-center items-center py-10">
      <ProgressSpinner style="width: 50px; height: 50px;" strokeWidth="5" animationDuration=".7s" />
    </div>

    <!-- Error Message -->
    <div v-else-if="errorMessage" class="p-error p-4 text-red-500">
      {{ errorMessage }}
    </div>

    <!-- Profile Details -->
    <div v-else class="profile-details relative">
      <div class="bg-white shadow-lg rounded-2xl flex flex-col items-center overflow-hidden">
        <!-- Vue Green Background -->
        <div class="w-full bg-gradient-to-r from-green-400 via-green-500 to-green-600  h-60 bg-[#42b883]"></div>

        <!-- Profile Image -->
        <div class="w-32 h-32 rounded-full bg-white border-4 border-green-400 -mt-16 flex items-center justify-center z-10">
          <img :src="user?.profile_photo_url || '/profile.svg'" alt="Profile Photo"
            class="w-full h-full rounded-full object-cover" />
        </div>

        <!-- User Info -->
        <div class="px-6 pt-6 pb-6 text-center">
          <h2 class="font-bold text-gray-800 text-2xl mb-2">{{ user.user_name }}</h2>
          <h4 class="text-gray-600 text-sm mb-2">
            {{ user.email_address }} · {{ user.phone_number }} · {{ user.location }}
          </h4>
        </div>
      </div>
    </div>

    <Divider class="my-6" />

    <!-- Two-Column Layout -->
    <div v-if="hasContent" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      
      <!-- Left Column: Tracks (2/3 width) -->
      <div class="lg:col-span-2 space-y-6">
        
        <!-- Tracks -->
        <section v-if="user.tracks.length">
          <h2 class="text-xl font-semibold mb-3 text-left">Tracks</h2>
          <div class="grid grid-cols-1 sm:grid-cols-1 md:grid-cols-1 lg:grid-cols-2 xl:grid-cols-2 gap-6 pb-4">
            <div v-for="t in user.tracks" :key="t.track_id" class="min-w-[280px] max-w-md">
              <TrackCard :track="t" />
            </div>
          </div>
        </section>

        <!-- Liked Events -->
        <section v-if="user.events.length" class="space-y-6">
          <h2 class="text-xl font-semibold mb-3 text-left">Liked Events</h2>
          <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-2 xl:grid-cols-2 gap-6 pb-4 max-h-[500px] overflow-y-auto">
            <EventCard v-for="event in user.events" :key="event.event_id" :event="event" :liked="true"
              @event-unliked="handleEventUnliked" @event-liked="handleEventLiked" />
          </div>
        </section>
      </div>

      <!-- Right Column: Pending Friend Requests & Friends -->
      <div class="lg:col-span-1 flex flex-col space-y-6 max-h-[900px]">
        
        <!-- Pending Friend Requests -->
        <section v-if="user.friendRequests?.users?.length" class="space-y-6  flex-grow lg:order-first">
          <h2 class="text-xl font-semibold mb-3 text-left overflow-y-auto">Pending Friend Requests</h2>
          <div class="space-y-1">
            <UserCard
              v-for="u in user.friendRequests.users"
              :key="u.user_id"
              :user="u"
              :accept="true"
              :reject="true"
              @accept="handleAccept(u.user_id)"
              @reject="handleReject(u.user_id)" />
          </div>
        </section>

        <section
        v-if="user.friends.length"
        class="space-y-6 flex-grow lg:order-first max-h-[500px] pr-2"
      >
        <h2 class="text-xl font-semibold mb-3 text-left">Friends · {{ user.friends.length }}</h2>
        <div class="space-y-4">
          <UserCard
            v-for="u in user.friends"
            :key="u.user_id"
            :user="u"
            :remove="true"
            @remove="handleRemove(u.user_id)"
          />
        </div>
      </section>


      </div>
    </div>
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
import {
  getFriends,
  getPendingFriendRequests,
  acceptFriendRequest,
  rejectFriendRequest, 
  removeFriend
} from "@/service/FriendService";


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
      return events.length || playlists.length || tracks.length;
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
            getFriends()
          ]);
          this.user = { ...u, events, playlists, tracks, friendRequests: null }; // No friendRequests for other users
        } else {
          // Fetch data for the logged-in user
          const [u, events, playlists, tracks, friends, friendRequests] = await Promise.all([
            UserService.getUser({ withCredentials: true }),
            EventService.getFavEventsForCurrentUser(),
            PlaylistService.getPlaylistsForUser(),
            getFavUserTracks(),
            getFriends(),
            getPendingFriendRequests(),
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
      acceptFriendRequest(userId)
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
      rejectFriendRequest(userId)
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
      removeFriend(userId)
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
