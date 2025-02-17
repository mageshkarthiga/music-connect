<!-- src/views/Home.vue -->
<template>
  <div class="p-d-flex p-flex-column" style="height: 100vh">
    <NavTabs :activeTab="activeTab" :tabs="tabs" @changeTab="changeTab" />
    <div class="p-flex-1 p-p-3">
      <component :is="activeTabComponent" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import NavTabs from "../components/NavTabs.vue";
import MusicTab from "../components/MusicTab.vue";
import EventsTab from "../components/EventsTab.vue";
import UsersTab from "../components/UsersTab.vue";

const activeTab = ref("music");

const tabs = [
  { id: "music", label: "Music" },
  { id: "events", label: "Events" },
  { id: "users", label: "Users" },
];

const changeTab = (tabId: string) => {
  activeTab.value = tabId;
};

const activeTabComponent = computed(() => {
  switch (activeTab.value) {
    case "music":
      return MusicTab;
    case "events":
      return EventsTab;
    case "users":
      return UsersTab;
    default:
      return MusicTab;
  }
});
</script>
