<script>
import { defineComponent, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import UserService from "@/service/UserService";
import { getAuth, onAuthStateChanged } from "firebase/auth";

export default defineComponent({
  name: "Profile",
  setup() {
    const placeholderImage = "https://picsum.photos/300/200";

    const route = useRoute();
    const user = ref(null);
    const loading = ref(true);

    const mockUser = {
      userName: "Mock User",
      emailAddress: "mock@example.com",
      phoneNumber: "123-456-7890",
      location: "Nowhere",
      profilePhotoUrl: "https://via.placeholder.com/120",

      events: [
        {
          event_id: 1,
          event_name: "TAEYEON CONCERT - The TENSE in SING",
          event_description: "TAEYEON TO RETURN TO SINGA",
          event_url: "https://ticketmaster.sg/activity/detail/25sg_taeyeon",
          event_image_url:
            "https://static.ticketmaster.sg/images/activity/25sg_taeyeon_0ac82c94b44eab77c29e2b00ec518b7c.jpg",
        },
        {
          event_id: 2,
          event_name: "RUSSELL PETERS RELAX WORLD TOUR",
          event_description: "Russell Peters' comedy tour.",
          event_url: "https://ticketmaster.sg/activity/2",
          event_image_url:
            "https://static.ticketmaster.sg/images/activity/25sg_russell_e9b8f4151423d68ecf2b41554373029e.jpg",
        },
        {
          event_id: 3,
          event_name: "BOYCE AVENUE LIVE IN SINGAPORE",
          event_description: "Boyce Avenue returns with acoustic performances.",
          event_url: "https://ticketmaster.sg/activity/3",
          event_image_url:
            "https://static.ticketmaster.sg/images/activity/25sg_boyceavenue_49a27ce4785567cffe6d0be75fd066bb.jpg",
        },
        {
          event_id: 4,
          event_name: "Wubai and China Blue Rock Star 2 世界巡回演唱会",
          event_description:
            "Wubai and China Blue’s world tour live in Singapore.",
          event_url: "https://ticketmaster.sg/activity/4",
          event_image_url:
            "https://static.ticketmaster.sg/images/activity/4.jpg",
        },
        {
          event_id: 5,
          event_name: "Jia Le Hokkien Hits Concert II",
          event_description:
            "A concert featuring popular Hokkien hits and artists.",
          event_url: "https://ticketmaster.sg/activity/5",
          event_image_url:
            "https://static.ticketmaster.sg/images/activity/5.jpg",
        },
        {
          event_id: 6,
          event_name: "Tchaikovsky Violin Concerto by SSO",
          event_description:
            "An evening of classical brilliance by the Singapore Symphony Orchestra.",
          event_url: "https://ticketmaster.sg/activity/6",
          event_image_url:
            "https://static.ticketmaster.sg/images/activity/6.jpg",
        },
        {
          event_id: 7,
          event_name: "Stand-Up Asia Comedy Fest",
          event_description:
            "Laugh-out-loud performances by top comedians from across Asia.",
          event_url: "https://ticketmaster.sg/activity/7",
          event_image_url:
            "https://static.ticketmaster.sg/images/activity/7.jpg",
        },
        {
          event_id: 8,
          event_name: "The Greatest Showman Live Tribute",
          event_description:
            "A musical tribute to the beloved film with live acts and sing-alongs.",
          event_url: "https://ticketmaster.sg/activity/8",
          event_image_url:
            "https://static.ticketmaster.sg/images/activity/8.jpg",
        },
        {
          event_id: 9,
          event_name: "Jazz in the Garden: Evening Sessions",
          event_description:
            "Smooth jazz vibes under the stars at Botanic Gardens.",
          event_url: "https://ticketmaster.sg/activity/9",
          event_image_url:
            "https://static.ticketmaster.sg/images/activity/9.jpg",
        },
        {
          event_id: 10,
          event_name: "BTS Fan Event – Lights & Memories",
          event_description:
            "A special event for ARMYs with BTS-themed activities and merch.",
          event_url: "https://ticketmaster.sg/activity/10",
          event_image_url:
            "https://static.ticketmaster.sg/images/activity/10.jpg",
        },
      ],
      tracks: [
        {
          track_id: 1,
          track_title: "number one girl",
          artist_id: 1,
          genre: "Pop",
          track_uri: "spotify:track:1lcBt7LoEikqYmhUoa2cez",
          track_spotify_id: "1lcBt7LoEikqYmhUoa2cez",
          image_url: placeholderImage,
        },
        {
          track_id: 2,
          track_title: "3am",
          artist_id: 1,
          genre: "R&B",
          track_uri: "spotify:track:3y4q6bBdbXsTlaPiwiiUfy",
          track_spotify_id: "3y4q6bBdbXsTlaPiwiiUfy",
          image_url: placeholderImage,
        },
        {
          track_id: 3,
          track_title: "two years",
          artist_id: 1,
          genre: "Indie",
          track_uri: "spotify:track:4HxGH28DitgAuuKpEVrLzN",
          track_spotify_id: "4HxGH28DitgAuuKpEVrLzN",
          image_url: placeholderImage,
        },
        {
          track_id: 4,
          track_title: "toxic till the end",
          artist_id: 1,
          genre: "Alternative",
          track_uri: "spotify:track:1z5ebC9238uGoBgzYyvGpQ",
          track_spotify_id: "1z5ebC9238uGoBgzYyvGpQ",
          image_url: placeholderImage,
        },
      ],

      playlists: [
        {
          playlist_id: 1,
          playlist_name: "My New Playlist",
          user_id: 76,
          image_url: placeholderImage,
        },
        {
          playlist_id: 2,
          playlist_name: "My Second Playlist",
          user_id: 76,
          image_url: placeholderImage,
        },
      ],

      artists: [
        {
          artist_id: 1,
          artist_name: "ROSE",
          artist_spotify_id: "3eVa5w3URK5duf6eyVDbu9",
          image_url: placeholderImage,
        },
        {
          artist_id: 2,
          artist_name: "Kendrick Lamar",
          artist_spotify_id: "2YZyLoL8N0Wb9xBt1NhZWg",
          image_url: placeholderImage,
        },
        {
          artist_id: 3,
          artist_name: "LE SSERAFIM",
          artist_spotify_id: "4SpbR6yFEvexJuaBpgAU5p",
          image_url: placeholderImage,
        },
        {
          artist_id: 4,
          artist_name: "Selena Gomez",
          artist_spotify_id: "0C8ZW7ezQVs4URX5aX7Kqx",
          image_url: placeholderImage,
        },
        {
          artist_id: 5,
          artist_name: "Bruno Mars",
          artist_spotify_id: "0du5cEVh5yTK9QJze8zA0C",
          image_url: placeholderImage,
        },
      ],
    };

    onMounted(async () => {
      const userId = Number(route.query.user_id);
      const auth = getAuth();

      if (!isNaN(userId)) {
        try {
          const fetchedUser = await UserService.getUser(userId, {
            withCredentials: true,
          });
          user.value = fetchedUser || mockUser;
        } catch (error) {
          console.error("Error fetching user:", error);
          user.value = mockUser;
        } finally {
          loading.value = false;
        }
      } else {
        onAuthStateChanged(auth, async (firebaseUser) => {
          if (firebaseUser) {
            try {
              const fetchedUser = await UserService.getUserByFirebaseUID(
                firebaseUser.uid
              );
              user.value = fetchedUser || mockUser;
              user.value = mockUser; //temp
            } catch (error) {
              console.error("Error fetching user by Firebase UID:", error);
              user.value = mockUser;
            } finally {
              loading.value = false;
            }
          } else {
            console.warn("No Firebase user signed in.");
            user.value = mockUser;
            loading.value = false;
          }
        });
      }
    });

    return {
      user,
      loading,
    };
  },
});
</script>
<template>
  <div
    class="profile-page"
    style="max-width: 1000px; margin: 2rem auto; text-align: center"
  >
    <div
      v-if="loading"
      class="loading"
      style="font-size: 1.2rem; padding: 2rem"
    >
      <p>Loading...</p>
    </div>

    <div v-else-if="user" class="profile-details p-card p-p-4 p-shadow-4">
      <img
        :src="user.profilePhotoUrl"
        alt="Profile Photo"
        class="profile-photo"
        style="
          width: 120px;
          height: 120px;
          object-fit: cover;
          border-radius: 50%;
          border: 3px solid var(--primary-color);
        "
      />
      <h1 class="p-mt-3">{{ user.userName }}</h1>
      <p><strong>Email:</strong> {{ user.emailAddress }}</p>
      <p><strong>Phone:</strong> {{ user.phoneNumber }}</p>
      <p><strong>Location:</strong> {{ user.location }}</p>
    </div>

    <div v-else class="error p-error" style="font-size: 1.2rem; padding: 2rem">
      <p>User not found.</p>
    </div>
  </div>

  <!-- Events Section -->
  <div v-if="user && user.events" class="p-4">
    <h2 class="text-xl font-semibold mb-3">Events</h2>
    <div class="flex space-x-4 overflow-x-auto pb-4">
      <div
        v-for="event in user.events"
        :key="event.event_id"
        class="min-w-[250px] max-w-[300px] flex-shrink-0"
      >
        <Card class="h-full">
          <template v-slot:title>
            <div class="flex items-center justify-between mb-0">
              <div class="font-semibold text-xl mb-4">
                {{ event.event_name }}
              </div>
              <Button icon="pi pi-calendar" class="p-button-text" />
            </div>
            <Menu id="event_menu" :popup="true" />
          </template>

          <template v-slot:content>
            <div class="flex justify-between flex-col">
              <img
                :src="event.event_image_url"
                alt="Event image"
                class="w-full h-40 object-cover rounded-xl mb-3"
              />
              <p class="text-sm text-gray-800 line-clamp-3 mt-1 mb-4">
                {{ event.event_description }}
              </p>
              <Button
                label="Buy Tickets"
                severity="success"
                class="w-full justify-center mt-auto"
                :pt="{ root: 'text-sm font-medium px-4 py-2 rounded-xl' }"
                :href="event.event_url"
                target="_blank"
              />
            </div>
          </template>
        </Card>
      </div>
    </div>
  </div>
  <!-- Tracks Section -->
  <div v-if="user && user.tracks" class="p-4">
    <h2 class="text-xl font-semibold mb-3">Tracks</h2>
    <div class="flex space-x-4 overflow-x-auto pb-4">
      <div
        v-for="track in user.tracks"
        :key="track.track_id"
        class="min-w-[250px] max-w-[300px] flex-shrink-0"
      >
        <Card>
          <template v-slot:title>
            <div class="flex items-center justify-between mb-0">
              <div class="font-semibold text-xl mb-4">
                {{ track.track_title }}
              </div>
              <Button icon="pi pi-music" class="p-button-text" />
            </div>
            <Menu id="track_menu" :popup="true" />
          </template>

          <template v-slot:content>
            <img
              :src="track.image_url"
              alt="Track image"
              class="w-full h-40 object-cover rounded-xl mb-3"
            />
            <div class="mb-2">
              <span class="text-sm font-semibold text-gray-800">Genre:</span>
              <p class="text-sm text-gray-800">{{ track.genre || "N/A" }}</p>
            </div>
            <Button
              label="Listen on Spotify"
              severity="success"
              class="w-full justify-center mt-auto"
              :pt="{ root: 'text-sm font-medium px-4 py-2 rounded-xl' }"
              :href="`https://open.spotify.com/track/${track.track_spotify_id}`"
              target="_blank"
            />
          </template>
        </Card>
      </div>
    </div>
  </div>
  <!-- Playlists Section -->
  <div v-if="user && user.playlists" class="p-4">
    <h2 class="text-xl font-semibold mb-3">Playlists</h2>
    <div class="flex space-x-4 overflow-x-auto pb-4">
      <div
        v-for="playlist in user.playlists"
        :key="playlist.playlist_id"
        class="min-w-[250px] max-w-[300px] flex-shrink-0"
      >
        <Card>
          <template v-slot:title>
            <div class="flex items-center justify-between mb-0">
              <div class="font-semibold text-xl mb-4">
                {{ playlist.playlist_name }}
              </div>
              <Button icon="pi pi-list" class="p-button-text" />
            </div>
            <Menu id="playlist_menu" :popup="true" />
          </template>

          <template v-slot:content>
            <img
              :src="playlist.image_url"
              alt="Playlist image"
              class="w-full h-40 object-cover rounded-xl mb-3"
            />
            <p class="text-sm text-gray-800 mt-auto">
              User ID: {{ playlist.user_id }}
            </p>
          </template>
        </Card>
      </div>
    </div>
  </div>
  <!-- Artists Section -->
  <div v-if="user && user.artists" class="p-4">
    <h2 class="text-xl font-semibold mb-3">Artists</h2>
    <div class="flex space-x-4 overflow-x-auto pb-4">
      <div
        v-for="artist in user.artists"
        :key="artist.artist_id"
        class="min-w-[250px] max-w-[300px] flex-shrink-0"
      >
        <Card>
          <template v-slot:title>
            <div class="flex items-center justify-between mb-0">
              <div class="font-semibold text-xl mb-4">
                {{ artist.artist_name }}
              </div>
              <Button icon="pi pi-user" class="p-button-text" />
            </div>
            <Menu id="artist_menu" :popup="true" />
          </template>

          <template v-slot:content>
            <img
              :src="artist.image_url"
              alt="Artist image"
              class="w-full h-40 object-cover rounded-xl mb-3"
            />
            <Button
              label="View on Spotify"
              severity="success"
              class="w-full justify-center mt-auto"
              :pt="{ root: 'text-sm font-medium px-4 py-2 rounded-xl' }"
              :href="`https://open.spotify.com/artist/${artist.artist_spotify_id}`"
              target="_blank"
            />
          </template>
        </Card>
      </div>
    </div>
  </div>
</template>
