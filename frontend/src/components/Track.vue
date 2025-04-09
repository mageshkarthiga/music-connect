<template>
  <div class="min-w-[250px] max-w-[300px] flex-shrink-0">
    <Card>
      <template v-slot:title>
        <div class="flex items-center justify-between mb-0">
          <div class="font-semibold text-xl mb-4">
            {{ track.track_title }}
          </div>
        </div>
      </template>
      <template v-slot:content>
        <img
          :src="track.track_image_url || fallbackImage"
          alt="Track image"
          class="w-full h-40 object-cover rounded-xl mb-3"
        />
        <p class="text-sm text-gray-700 line-clamp-1">
          {{ track.genre }}
        </p>
      </template>
    </Card>
  </div>
</template>

<script>
export default {
  props: {
    track: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      fallbackImage: "https://via.placeholder.com/300x160?text=Track+Image",
    };
  },
};
</script>
