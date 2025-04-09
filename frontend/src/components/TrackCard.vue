<template>
  <div
    class="flex items-center gap-4 bg-zinc-900 rounded-lg p-3 min-w-[280px] max-w-md hover:bg-zinc-800 transition cursor-pointer"
    @click="openTrack"
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
      fallbackImage: "https://picsum.photos/300/200",
    };
  },
  methods: {
    openTrack() {
      if (this.track?.track_uri) {
        window.open(this.track.track_uri, "_blank");
      }
    },
  },
};
</script>
