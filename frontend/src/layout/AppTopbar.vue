<script setup>
import {
  useLayout,
  toggleMenu,
  toggleDarkMode,
} from "@/layout/composables/stateConfig";
import AppConfigurator from "./AppConfigurator.vue";
import { useRoute, useRouter } from "vue-router";
import { signOut } from "firebase/auth";
import { auth } from "@/firebase/firebase"; // adjust path if needed

const { isDarkTheme } = useLayout();
const router = useRouter();
function goToProfile() {
  router.push("/profile");
}
function goToSearchBar() {
  router.push("/pages/search");
}
function goToChat() {
  router.push("/pages/chat"); 
}
function goToMap() {
  router.push("/uikit/map");
}
function logout() {
  signOut(auth)
    .then(() => {
      router.replace("/auth/login");
    })
    .catch((error) => {
      console.error("Logout failed:", error);
    });
}
</script>

<template>
  <div class="layout-topbar">
    <div class="layout-topbar-logo-container">
      <button
        class="layout-menu-button layout-topbar-action"
        @click="toggleMenu"
      >
        <i class="pi pi-bars"></i>
      </button>
      <router-link to="/pages/home" class="layout-topbar-logo">
        <Avatar image="/public/logo.png" size="large" />
        <span class="font-bold text-xl">Music Connect</span>
      </router-link>
    </div>

    <div class="layout-topbar-actions">
      <div class="layout-config-menu">
        <button
          type="button"
          class="layout-topbar-action"
          @click="toggleDarkMode"
        >
          <i
            :class="['pi', { 'pi-moon': isDarkTheme, 'pi-sun': !isDarkTheme }]"
          ></i>
        </button>
      </div>

      <button
        class="layout-topbar-menu-button layout-topbar-action"
        v-styleclass="{
          selector: '@next',
          enterFromClass: 'hidden',
          enterActiveClass: 'animate-scalein',
          leaveToClass: 'hidden',
          leaveActiveClass: 'animate-fadeout',
          hideOnOutsideClick: true,
        }"
      >
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
          <button
            type="button"
            class="layout-topbar-action"
            @click="goToProfile"
          >
            <i class="pi pi-user"></i>
            <span>Profile</span>
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
