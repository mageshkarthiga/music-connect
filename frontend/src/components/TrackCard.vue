<template>
  <div
    class="dark flex items-center gap-4 bg-surface-500 rounded-lg p-3 min-w-[280px] max-w-md hover:bg-surface-400 dark:hover:bg-surface-600 transition cursor-pointer"
    @click="selectTrack"
  >
    <img
      :src="track.track_image_url || fallbackImage"
      alt="Track cover"
      class="w-12 h-12 object-cover rounded-md"
    />
    <div class="text-white font-medium text-sm truncate">
      {{ track.track_title }}
    </div>
  </div>
</template>

<script>
export default {
  props: {
    track: { type: Object, required: true },
  },
  data() {
    return {
      fallbackImage: "https://picsum.photos/300/200", // Fallback image if no track image is provided
    };
  },
  methods: {
    selectTrack() {
      // Emit an event with the track's URI when the card is clicked
      this.$emit("trackSelected", this.track.track_uri);
    },
  },
};
</script>