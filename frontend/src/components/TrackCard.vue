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

    <div
      v-if="state === 'redirect'"
      @click.stop="toggleLike"
      class="absolute right-2 top-1/2 transform -translate-y-1/2"
    >
      <i
        class="pi"
        :class="likedStatus ? 'pi-heart-fill text-red-500' : 'pi-heart'"
        style="font-size: 1.2rem;"
      ></i>
    </div>
  </div>
</template>

<script>
import { API_BASE_URL } from "@/service/apiConfig";
import axios from "axios";

export default {
  props: {
    track: { type: Object, required: true },
    state: { type: String, required: true },
    liked: { type: Boolean, default: false }, // Default to true
  },
  data() {
    return {
      fallbackImage: "https://picsum.photos/300/200",
      likedStatus: this.liked, // Bind the liked state to a data property
    };
  },
  methods: {
    toggleLike() {
      const trackId = this.track.track_id;
      const updateStatus = !this.likedStatus;

      // Toggle like/unlike logic
      if (updateStatus) {
        // If liking the track
        axios.put(`${API_BASE_URL}//likeTrack/${trackId}`, { is_liked: true }, {
          withCredentials: true,
        });
        this.$emit("track-liked", trackId);
      } else {
        // If unliking the track
        axios.put(`${API_BASE_URL}//unlikeTrack/${trackId}`, { is_liked: false }, {
          withCredentials: true,
        });
        this.$emit("track-unliked", trackId);
      }

      // Update the liked status
      this.likedStatus = updateStatus;

      // Emit the like/unlike event
      this.$emit(updateStatus ? "track-liked" : "track-unliked", trackId);
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
