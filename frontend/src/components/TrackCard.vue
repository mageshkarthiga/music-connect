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
/* Track card styles */
.track-card {
  padding: 16px;
  border-radius: 12px;
  border: 1px solid rgba(184, 184, 184, 0.5); /* Light border for light mode */
  cursor: pointer;
  transition: 0.2s;
}

/* Hover effect for light mode */
.bg-surface-400:hover {
  background-color: rgba(245, 245, 245, 0.5);
}

/* Dark mode background and border */
.dark .track-card {
  background-color: #101318; /* Darker background for track cards in dark mode */
  border: 1px solid #4a5568; /* Darker border for dark mode */
}

/* Light mode specific card styles */
.track-card.light {
  background-color: #ffffff; /* Set light mode card background to #ffffff */
  border: 1px solid #ddd; /* Light border */
}

/* Hover effect for dark mode */
.dark .bg-surface-400:hover {
  background-color: rgba(184, 184, 184, 0.5) /* Darker hover effect for dark mode */
}

/* Highlight the selected card */
.track-card.selected {
  border: 2px solid #585858;
  background-color: #f0f9ff;
}

/* Image style */
.track-image {
  width: 100%;
  height: auto;
  border-radius: 8px;
}

/* Track info style */
.track-info {
  margin-top: 0.5rem;
}
</style>
