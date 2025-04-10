<!-- <script setup>
import { ref } from "vue";

import AppMenuItem from "./AppMenuItem.vue";

const model = ref([
  {
    label: "Home",
    items: [{ label: "Dashboard", icon: "pi pi-fw pi-home", to: "/" }],
  },
  {
    label: "UI Components",
    items: [
      {
        label: "Form Layout",
        icon: "pi pi-fw pi-id-card",
        to: "/uikit/formlayout",
      },
      { label: "Input", icon: "pi pi-fw pi-check-square", to: "/uikit/input" },
      {
        label: "Button",
        icon: "pi pi-fw pi-mobile",
        to: "/uikit/button",
        class: "rotated-icon",
      },
      { label: "Table", icon: "pi pi-fw pi-table", to: "/uikit/table" },
      { label: "List", icon: "pi pi-fw pi-list", to: "/uikit/list" },
      { label: "Tree", icon: "pi pi-fw pi-share-alt", to: "/uikit/tree" },
      { label: "Panel", icon: "pi pi-fw pi-tablet", to: "/uikit/panel" },
      { label: "Overlay", icon: "pi pi-fw pi-clone", to: "/uikit/overlay" },
      { label: "Media", icon: "pi pi-fw pi-image", to: "/uikit/media" },
      { label: "Menu", icon: "pi pi-fw pi-bars", to: "/uikit/menu" },
      { label: "Message", icon: "pi pi-fw pi-comment", to: "/uikit/message" },
      { label: "File", icon: "pi pi-fw pi-file", to: "/uikit/file" },
      { label: "Chart", icon: "pi pi-fw pi-chart-bar", to: "/uikit/charts" },
      {
        label: "Timeline",
        icon: "pi pi-fw pi-calendar",
        to: "/uikit/timeline",
      },
      { label: "Misc", icon: "pi pi-fw pi-circle", to: "/uikit/misc" },
      { label: "Map", icon: "pi pi-fw pi-map", to: "/uikit/map" },
    ],
  },
  {
    label: "Pages",
    icon: "pi pi-fw pi-briefcase",
    to: "/pages",
    items: [
      {
        label: "Landing",
        icon: "pi pi-fw pi-globe",
        to: "/landing",
      },
      {
        label: "Auth",
        icon: "pi pi-fw pi-user",
        items: [
          {
            label: "Login",
            icon: "pi pi-fw pi-sign-in",
            to: "/auth/login",
          },
          {
            label: "Error",
            icon: "pi pi-fw pi-times-circle",
            to: "/auth/error",
          },
          {
            label: "Access Denied",
            icon: "pi pi-fw pi-lock",
            to: "/auth/access",
          },
        ],
      },
      {
        label: "Crud",
        icon: "pi pi-fw pi-pencil",
        to: "/pages/crud",
      },
      {
        label: "Not Found",
        icon: "pi pi-fw pi-exclamation-circle",
        to: "/pages/notfound",
      },
      {
        label: "Create Account",
        icon: "pi pi-fw pi-user-plus",
        to: "/createaccount",
      },
      {
        label: "Profile",
        icon: "pi pi-fw pi-user",
        to: "/profile",
      },
    ],
  },
  {
    label: "Hierarchy",
    items: [
      {
        label: "Submenu 1",
        icon: "pi pi-fw pi-bookmark",
        items: [
          {
            label: "Submenu 1.1",
            icon: "pi pi-fw pi-bookmark",
            items: [
              { label: "Submenu 1.1.1", icon: "pi pi-fw pi-bookmark" },
              { label: "Submenu 1.1.2", icon: "pi pi-fw pi-bookmark" },
              { label: "Submenu 1.1.3", icon: "pi pi-fw pi-bookmark" },
            ],
          },
          {
            label: "Submenu 1.2",
            icon: "pi pi-fw pi-bookmark",
            items: [{ label: "Submenu 1.2.1", icon: "pi pi-fw pi-bookmark" }],
          },
        ],
      },
      {
        label: "Submenu 2",
        icon: "pi pi-fw pi-bookmark",
        items: [
          {
            label: "Submenu 2.1",
            icon: "pi pi-fw pi-bookmark",
            items: [
              { label: "Submenu 2.1.1", icon: "pi pi-fw pi-bookmark" },
              { label: "Submenu 2.1.2", icon: "pi pi-fw pi-bookmark" },
            ],
          },
          {
            label: "Submenu 2.2",
            icon: "pi pi-fw pi-bookmark",
            items: [{ label: "Submenu 2.2.1", icon: "pi pi-fw pi-bookmark" }],
          },
        ],
      },
    ],
  },
  {
    label: "Get Started",
    items: [
      {
        label: "Documentation",
        icon: "pi pi-fw pi-book",
        to: "/documentation",
      },
      {
        label: "View Source",
        icon: "pi pi-fw pi-github",
        url: "https://github.com/primefaces/sakai-vue",
        target: "_blank",
      },
    ],
  },
]);
</script>

<template>
  <ul class="layout-menu">
    <template v-for="(item, i) in model" :key="item">
      <app-menu-item
        v-if="!item.separator"
        :item="item"
        :index="i"
      ></app-menu-item>
      <li v-if="item.separator" class="menu-separator"></li>
    </template>
  </ul>
</template>

<style lang="scss" scoped></style> -->

<template>
  <!-- Sidebar -->
  <div class="sidebar">
    <!-- Library Header -->
    <div class="library-header">
      <h4>Your Library</h4>
      <Button @click="openAddPlaylistDialog" class="add-playlist-btn">
        <i class="pi pi-plus"></i> 
      </Button>
    </div>

    <!-- Playlists List -->
    <div v-if="user.playlists.length" class="playlist-list">
      <div
        v-for="playlist in user.playlists"
        :key="playlist.playlist_id"
        class="playlist-item"  
      >
        <img :src="playlist.playlist_image_url" alt="Playlist Image" class="playlist-image" />
        <div class="playlist-info">
          <span class="playlist-name">{{ playlist.playlist_name }}</span>
          <span class="playlist-username"> Playlist - {{ currentUser?.user_name }}</span>
        </div>
      </div>
    </div>

    <div v-else>
      <p>No playlists available.</p>
    </div>

    <!-- Playlist Dialog -->
    <AddPlaylistDialog 
      v-if="showPlaylistDialog" 
      @close="closeAddPlaylistDialog" 
      @playlist-added="addNewPlaylist"
    />
  </div>
</template>

<script>
import axios from "axios";
import { API_BASE_URL } from "@/service/apiConfig";
import AddPlaylistDialog from "@/components/AddPlaylistDialog.vue";

export default {
  components: {
    AddPlaylistDialog,
  },
  data() {
    return {
      user: { playlists: [] },
      currentUser: null,
      errorMessage: "",
      darkTheme: localStorage.getItem("theme") === "dark",
      showPlaylistDialog: false,
    };
  },
  methods: {
    async getCurrentUser() {
      try {
        const response = await axios.get(`${API_BASE_URL}/me`, {
          withCredentials: true,
        });
        this.currentUser = response.data;
      } catch (err) {
        this.handleError(err, "current user");
      }
    },

    async getPlaylistsForUser() {
      try {
        const response = await axios.get(`${API_BASE_URL}/me/playlists`, {
          withCredentials: true,
        });

        if (Array.isArray(response.data)) {
          this.user.playlists = response.data;
        } else {
          throw new Error("Invalid playlists data format");
        }
      } catch (err) {
        this.handleError(err, "playlists");
      }
    },

    handleError(error, dataType) {
      console.error(`${dataType} fetch error:`, error);
      this.errorMessage = error.response?.data?.message || `Failed to fetch ${dataType}.`;
    },

    openAddPlaylistDialog() {
      this.showPlaylistDialog = true;
    },

    closeAddPlaylistDialog() {
      this.showPlaylistDialog = false;
    },

    addNewPlaylist(newPlaylist) {
      this.user.playlists.push(newPlaylist);
      this.closeAddPlaylistDialog();
    },

    toggleTheme() {
      this.darkTheme = !this.darkTheme;
      localStorage.setItem("theme", this.darkTheme ? "dark" : "light");
    },
  },
  mounted() {
    this.getCurrentUser();
    this.getPlaylistsForUser();
  },
};
</script>



<style lang="scss" scoped>
.sidebar {
  padding: 1rem;
  width: 250px;
  height: 100%;
  display: flex;
  flex-direction: column;
  transition: background-color 0.3s ease, color 0.3s ease;
}

.dark {
  background-color: #121212;
  color: white;
}

.light {
  background-color: #f4f4f4;
  color: black;
}

h4 {
  color: inherit;
}

.library-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.add-playlist-btn {
  background-color: #1db954;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 50%;
  cursor: pointer;
  font-size: 1.5rem;
  width: 40px;
  height: 40px;
}

.add-playlist-btn:hover {
  background-color: #1ed760;
}

.theme-toggle-btn {
  background-color: transparent;
  border: none;
  color: inherit;
  font-size: 1.5rem;
  cursor: pointer;
}

.playlist-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.playlist-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.5rem;
  border-radius: 8px;
  transition: background-color 0.3s ease;
}

.dark .playlist-item {
  background-color: #282828;
}

.light .playlist-item {
  background-color: #ffffff;
}

.playlist-item:hover {
  background-color: rgba(125, 125, 125, 0.1)
}
.dark .playlist-item:hover {
  background-color: #3a3a3a;
}

.light .playlist-item:hover {
  background-color: #cacaca;
}

.playlist-image {
  width: 50px;
  height: 50px;
  border-radius: 8px;
  object-fit: cover;
}

.playlist-info {
  display: flex;
  flex-direction: column;
}

.playlist-name {
  font-size: 1rem;
  font-weight: 600;
}

.playlist-username {
  font-size: 0.875rem;
  color: rgba(125, 125, 125, 0.9)
}

.dark .playlist-name {
  color: rgba(255, 255, 255, 0.7)
}
.light .playlist-username {
  color: rgb(0, 0, 0, 0.7)
}
</style>
