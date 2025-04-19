<script>
import { useLayout, toggleMenu, toggleDarkMode } from "@/layout/composables/stateConfig";
import { signOut } from "firebase/auth";
import { auth } from "@/firebase/firebase"; 
import { getPendingFriendRequests } from "@/service/FriendService";

export default {
  name: "AppTopbar",
  data() {
    return {
      isDarkTheme: useLayout().isDarkTheme,
      pendingFriendRequestsCount: 0,
    };
  },
  mounted() {
    this.fetchPendingFriendRequests();
    this.refreshInterval = setInterval(() => {
      this.fetchPendingFriendRequests();
    }, 60000); 
  },
  beforeDestroy() {
    clearInterval(this.refreshInterval); 
  },
  methods: {
    toggleMenu,
    toggleDarkMode,
    goToProfile() {
      this.$router.push("/profile");
    },
    goToSearchBar() {
      this.$router.push("/pages/search");
    },
    goToChat() {
      this.$router.push("/pages/chat");
    },
    goToMap() {
      this.$router.push("/pages/map");
    },
    logout() {
      signOut(auth)
        .then(() => {
          this.$router.replace("/auth/login");
        })
        .catch((error) => {
          console.error("Logout failed:", error);
        });
    },
    async fetchPendingFriendRequests() {
      try {
        const response = await getPendingFriendRequests();
        this.pendingFriendRequestsCount = response.count;
      } catch (error) {
        console.error("Error fetching pending friend requests:", error);
      }
    },
  },
};
</script>

<template>
  <div class="layout-topbar">
    <div class="layout-topbar-logo-container">
      <button class="layout-menu-button layout-topbar-action" @click="toggleMenu">
        <i class="pi pi-bars"></i>
      </button>
      <router-link to="/pages/home" class="layout-topbar-logo">
        <Avatar image="/public/logo.png" size="large" />
        <span class="font-bold text-xl">Music Connect</span>
      </router-link>
    </div>

    <div class="layout-topbar-actions">
      <div class="layout-config-menu">
        <button type="button" class="layout-topbar-action" @click="toggleDarkMode">
          <i :class="['pi', { 'pi-moon': isDarkTheme, 'pi-sun': !isDarkTheme }]"></i>
        </button>
      </div>

      <button class="layout-topbar-menu-button layout-topbar-action" v-styleclass="{
        selector: '@next',
        enterFromClass: 'hidden',
        enterActiveClass: 'animate-scalein',
        leaveToClass: 'hidden',
        leaveActiveClass: 'animate-fadeout',
        hideOnOutsideClick: true,
      }">
        <i class="pi pi-ellipsis-v"></i>
      </button>

      <div class="layout-topbar-menu hidden lg:block">
        <div class="layout-topbar-menu-content">
          <button type="button" class="layout-topbar-action" @click="goToMap">
            <i class="pi pi-map"></i>
            <span>Map</span>
          </button>
          <button type="button" class="layout-topbar-action" @click="goToSearchBar">
            <i class="pi pi pi-search"></i>
            <span>Search</span>
          </button>
          <button type="button" class="layout-topbar-action" @click="goToChat">
            <i class="pi pi-comments"></i>
            <span>Chat</span>
          </button>
          <button type="button" class="profile-button" @click="goToProfile">
            <div class="icon-wrapper">
              <i class="pi pi-user user-icon"></i>
              <Badge v-if="pendingFriendRequestsCount > 0" :value="pendingFriendRequestsCount" severity="danger"
                class="badge-overlay" />
            </div>
          </button>
          <button type="button" class="layout-topbar-action" @click="logout">
            <i class="pi pi-sign-out"></i>
            <span>Logout</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-button {
  background: none;
  border: none;
  padding: 0;
  position: relative;
}

.icon-wrapper {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.3s ease;
}

.icon-wrapper:hover {
  background-color: var(--surface-hover);
}

.user-icon {
  font-size: 1.4rem;
}

.badge-overlay {
  position: absolute;
  top: 5px;
  right: 6px;
  transform: translate(50%, -50%);
  font-size: 0.8rem;
  padding: 0.1rem 0.2rem;
  line-height: 1;
  border-radius: 9999px;
}
</style>