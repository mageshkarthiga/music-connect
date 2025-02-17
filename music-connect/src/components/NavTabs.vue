<!-- src/components/NavTabs.vue -->
<template>
  <TabMenu :model="tabItems" :activeIndex="activeIndex" @tab-change="onTabChange" />
</template>

<script setup lang="ts">
import { defineProps, computed, defineEmits } from "vue";
import TabMenu from "primevue/tabmenu";

const props = defineProps<{
  activeTab: string;
  tabs: { id: string; label: string }[];
}>();

const emit = defineEmits<{
  (e: "changeTab", tabId: string): void;
}>();

const activeIndex = computed(() => {
  return props.tabs.findIndex((tab) => tab.id === props.activeTab);
});

const tabItems = computed(() => {
  return props.tabs.map((tab) => ({
    label: tab.label,
  }));
});

const onTabChange = (event: { index: number }) => {
  const tabId = props.tabs[event.index].id;
  emit("changeTab", tabId);
};
</script>
