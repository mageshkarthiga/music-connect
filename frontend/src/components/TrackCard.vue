<template>
  <div
    class="dark flex items-center gap-4 bg-surface-500 rounded-lg p-3 min-w-[280px] max-w-md hover:bg-surface-400 dark:hover:bg-surface-600 transition cursor-pointer"
    @click="handleClick"
  >
    <img
      :src="track.track_image_url || fallbackImage"
      alt="Track cover"
      class="w-12 h-12 object-cover rounded-md"
    />
    <div class="text-white font-medium text-sm truncate">
      {{ track.track_title }}
    </div>
    <!-- Display selected status if in select state -->
    <div v-if="state === 'select'" v-show="isSelected" class="text-sm text-blue-500">
      Selected
    </div>
  </div>
</template>

<script>
export default {
  props: {
    track: { type: Object, required: true },
    state: { type: String, required: true }, // The state to control behavior ('redirect' or 'select')
    selectedTracks: { type: Array, required: true },
  },
  data() {
    return {
      fallbackImage: "https://picsum.photos/300/200",
    };
  },
  computed: {
    isSelected() {
      return this.selectedTracks.includes(this.track.track_id);
    },
  },
  methods: {
    handleClick() {
      if (this.state === 'redirect') {
        // If state is "redirect", open the track's URI
        if (this.track?.track_uri) {
          window.open(this.track.track_uri, "_blank");
        }
      } else if (this.state === 'select') {
        // If state is "select", toggle the track selection
        this.toggleTrackSelection();
      }
    },
    toggleTrackSelection() {
      const index = this.selectedTracks.indexOf(this.track.track_id);
      if (index > -1) {
        // Remove from selected
        this.selectedTracks.splice(index, 1);
      } else {
        // Add to selected
        this.selectedTracks.push(this.track.track_id);
      }
    },
  },
};
</script>

<style scoped>
.dark {
  color: #fff;
}

.track-card {
  padding: 12px;
  border-radius: 12px;
  border: 1px solid #ddd;
  cursor: pointer;
  transition: 0.2s;
}

.track-card.selected {
  border: 2px solid #3b82f6;
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
