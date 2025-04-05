// src/stores/spotifyStore.js
import { defineStore } from 'pinia';

export const useSpotifyStore = defineStore('spotify', {
  state: () => ({
    currentUri: null,
    // Other state properties related to Spotify
  }),
  actions: {
    setCurrentUri(uri) {
      this.currentUri = uri;
    },
    // Other actions related to Spotify
  },
});
