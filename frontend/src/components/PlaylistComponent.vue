<template>
    <div class="min-w-[250px] max-w-[300px] flex-shrink-0">
      <Card>
        <template v-slot:title>
          <div class="flex items-center justify-between mb-0">
            <div class="font-semibold text-xl mb-4">{{ playlist.playlist_name }}</div>
          </div>
        </template>
        <template v-slot:content>
          <img
            :src="playlist.image_url"
            alt="Playlist image"
            class="w-full h-40 object-cover rounded-xl mb-3"
          />
        </template>
      </Card>
    </div>
  </template>
  
  <script>
  export default {
    props: {
      playlist: Object, // Playlist data passed as prop
    },
  };
  </script>
  