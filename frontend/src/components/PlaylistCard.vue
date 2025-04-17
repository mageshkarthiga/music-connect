<template>
    <div class="min-w-[250px] max-w-[300px] flex-shrink-0">
      <Card @click="navigateToPlaylist" class="cursor-pointer">
        <template v-slot:title>
          <div class="flex items-center justify-between mb-0">
            <div class="font-semibold text-xl mb-4">
              {{ playlist.playlist_name }}
            </div>
          </div>
        </template>
        <template v-slot:content>
          <img :src="playlist.playlist_image_url" alt="Playlist image" class="w-full h-40 object-cover rounded-xl mb-3" />
        </template>
      </Card>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      playlist: { type: Object, required: true },
    },
    data() {
      return {
        fallbackImage: "https://picsum.photos/300/200", // Fallback image if no cover image is provided
      };
    },
    methods: {
      navigateToPlaylist() {
        this.$router.push({
          name: "playlist",
          params:
          {
            playlist_id: this.playlist.playlist_id,
            playlist_name: this.playlist.playlist_name
          }
        });
      },
    },
  };
  </script>
  