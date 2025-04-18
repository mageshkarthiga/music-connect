<template>
  <div
    class="track-card flex items-center gap-4 rounded-lg p-8 min-w-[280px] max-w-md hover:bg-surface-400 dark:bg-surface-900 dark:hover:bg-surface-800 transition cursor-pointer light:bg-surface-800 light:hover:bg-surface-300 relative"
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

    <!-- Selected status in select state -->
    <div v-if="state === 'select'" v-show="isSelected" class="text-sm text-blue-500">
      Selected
    </div>

    <div
  v-if="state === 'redirect'"
  @click.stop="toggleLike"
  class="absolute right-2 top-1/2 transform -translate-y-1/2"
>
  <!-- <i
    class="pi"
    :class="liked ? 'pi-heart-fill text-red-500' : 'pi-heart'"
    style="font-size: 1.2rem;"
  ></i> -->
</div>
  </div>
</template>

<script>
export default {
  props: {
    track: { type: Object, required: true },
    state: { type: String, required: true }, // 'redirect' or 'select'
    selectedTracks: { type: Array, required: true },
    liked: { type: Boolean, default: false }, // whether the track is liked
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
      if (this.state === "redirect") {
        this.$emit("track-selected", this.track.track_uri);
      } else if (this.state === "select") {
        this.toggleTrackSelection();
      }
    },
    toggleTrackSelection() {
      const index = this.selectedTracks.indexOf(this.track.track_id);
      if (index > -1) {
        this.selectedTracks.splice(index, 1);
      } else {
        this.selectedTracks.push(this.track.track_id);
      }
    },
    toggleLike() {
      if (this.liked) {
        this.$emit("track-unliked", this.track.track_id);
      } else {
        this.$emit("track-liked", this.track.track_id);
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
