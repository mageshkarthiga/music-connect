<template>
  <div
    class="track-card flex items-center gap-4 rounded-lg p-8 min-w-[280px] max-w-md hover:bg-surface-400 dark:bg-surface-900 dark:hover:bg-surface-800 transition cursor-pointer light:bg-surface-800 light:hover:bg-surface-300"
    @click="handleClick"
  >
    <img
      :src="track.track_image_url || fallbackImage"
      alt="Track cover"
      class="w-12 h-12 object-cover rounded-md"
    />
    <div class="text ml-4">
      {{ track.track_title }}
    </div>
    <!-- Display selected status if in select state -->
    <div v-if="state === 'select' && isSelected" class="text-sm text-blue-500">
      Selected
    </div>
  </div>
</template>

<script>
export default {
  props: {
    track: { type: Object, required: true },
    state: { type: String, required: true }, // 'redirect' or 'select'
    selectedTracks: { type: Array, required: true },
  },
  data() {
    return {
      fallbackImage: "https://picsum.photos/300/200",
    };
  },
  computed: {
    isSelected() {
      return Array.isArray(this.selectedTracks) && this.selectedTracks.includes(this.track.track_id);
    },
  },
  methods: {
    handleClick() {
      if (this.state === 'redirect') {
        if (this.track?.track_uri) {
          window.open(this.track.track_uri, "_blank");
        }
      } else if (this.state === 'select') {
        this.$emit('toggle', this.track.track_id);
      }
    },
  },
};
</script>

<style scoped>
.track-card {
  padding: 16px;
  border-radius: 12px;
  border: 1px solid rgba(184, 184, 184, 0.5);
  cursor: pointer;
  transition: 0.2s;
}

.bg-surface-400:hover {
  background-color: rgba(245, 245, 245, 0.5);
}

.dark .track-card {
  background-color: #101318;
  border: 1px solid #4a5568;
}

.track-card.light {
  background-color: #ffffff;
  border: 1px solid #ddd;
}

.dark .bg-surface-400:hover {
  background-color: rgba(184, 184, 184, 0.5);
}

.track-card.selected {
  border: 2px solid #585858;
  background-color: #f0f9ff;
}

.track-image {
  width: 100%;
  height: auto;
  border-radius: 8px;
}

.track-info {
  margin-top: 0.5rem;
}
</style>