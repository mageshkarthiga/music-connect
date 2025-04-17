<template>
    <div class="playlist-page">
      <div v-if="playlist?.tracks?.length">
        <TrackList :tracks="playlist.tracks" />
      </div>
  
      <div v-else-if="playlist" class="no-tracks">
        <p>This playlist is empty.</p>
      </div>
  
      <div v-else class="loading">
        <p>Loading playlist...</p>
      </div>
    </div>
  </template>
  
  <script>
  import axios from "axios";
  import { API_BASE_URL } from "@/service/apiConfig";
  import TrackList from "@/components/TrackList.vue";
  import PlaylistService from "@/service/PlaylistService";
  export default {
    name: "Playlist",
    components: {
      TrackList,
    },
    data() {
      return {
        playlist: null,
      };
    },
    async mounted() {
  const playlistId = this.$route.query.id;
  console.log("Playlist ID:", playlistId); 

  try {
    this.playlist = await PlaylistService.getPlaylistById(playlistId);
    console.log("Playlist data:", this.playlist); 
  } catch (error) {
    console.error("Error fetching playlist:", error);
  }
},

  };
  </script>
  
  <style scoped>
  .playlist-page {
    padding: 2rem;
  }
  
  .loading,
  .no-tracks {
    text-align: center;
    margin-top: 2rem;
    color: #aaa;
  }
  </style>
  