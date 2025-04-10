<template>
    <div>
      <label for="searchTracks">Search Tracks</label>
      <input
        v-model="searchQuery"
        id="searchTracks"
        type="text"
        placeholder="Search for tracks..."
        @input="filterTracks"
      />
    
      <ul>
        <li v-for="track in filteredTracks" :key="track.TrackID">
          <input type="checkbox" :value="track.TrackID" v-model="selectedTracks" />
          {{ track.TrackTitle }} - {{ track.Artists.map(artist => artist.ArtistName).join(', ') }}
        </li>
      </ul>
    </div>
  </template>
    
  <script>
  import { ref, computed, watch } from "vue";
  import axios from "axios";
  import { API_BASE_URL } from "@/service/apiConfig";
  
  export default {
    props: {
      selectedTracks: {
        type: Array,
        required: true,
      },
    },
    setup(props) {
      const tracks = ref([]);
      const searchQuery = ref("");
      const selectedTracks = props.selectedTracks;
    
      // Fetch all tracks on mount
      const fetchTracks = async () => {
        try {
          const response = await axios.get(`${API_BASE_URL}/tracks`, {
            withCredentials: true,
          });
          tracks.value = response.data;
        } catch (error) {
          console.error("Error fetching tracks:", error);
        }
      };
    
      // Filter tracks based on search query
      const filterTracks = () => {
        // If there's a search query, filter the tracks
        if (searchQuery.value) {
          return tracks.value.filter(track =>
            track.TrackTitle && track.TrackTitle.includes(searchQuery.value)
          );
        }
        // Return all tracks if no search query
        return tracks.value;
      };
    
      const filteredTracks = computed(filterTracks);
    
      // Fetch tracks when the component is mounted
      fetchTracks();
    
      return {
        tracks,
        searchQuery,
        selectedTracks,
        filteredTracks,
      };
    },
  };
  </script>
  