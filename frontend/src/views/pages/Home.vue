<template>
    <div class="home-page">
      <!-- Grid Section: Events and Music -->
      <div class="grid">
        <div class="row">
          <!-- Events Section -->
          <div class="col-6">
            <div class="row">
              <!-- Event Cards -->
              <div class="col-4 col-md-4">
                <Card>
                  <h3>Event 1: Concert</h3>
                  <p>Date: 10th April 2025</p>
                  <Button>RSVP</Button>

                </Card>
              </div>
  
              <div class="col-4 col-md-4">
                <Card>
                  <h3>Event 2: Festival</h3>
                  <p>Date: 20th May 2025</p>
                  <Button>RSVP</Button>
                </Card>
              </div>
  
              <div class="col-4 col-md-4">
                <Card>
                  <h3>Event 3: Festival</h3>
                  <p>Date: 20th May 2025</p>
                  <Button>RSVP</Button>
                </Card>
              </div>
            </div>
          </div>
  
          <!-- Music Section -->
          <div class="col-6">
            <div class="row">
              <!-- Music Cards -->
              <div class="col-4 col-md-4" v-for="track in tracks" :key="track.id">
                <Card>
                  <h3>{{ track.name }}</h3>
                  <p>Artist: {{ track.artist }}</p>
                  <Button @click="playTrack(track.uri)">Play</Button>
                </Card>
              </div>
            </div>
          </div>
        </div>
      </div>
  
      <!-- Fixed Spotify Player at Bottom -->
      <!-- <SpotifyPlayer /> -->
    </div>
  </template>
  
  <script>

  
  export default {

    data() {
      return {
        tracks: [
          {
            id: 1,
            name: "Song 1",
            artist: "Artist A",
            uri: "spotify:track:1",
          },
          {
            id: 2,
            name: "Song 2",
            artist: "Artist B",
            uri: "spotify:track:2",
          },
          {
            id: 3,
            name: "Song 3",
            artist: "Artist C",
            uri: "spotify:track:3",
          },
        ],
      };
    },
    methods: {
      playTrack(uri) {
        console.log("Playing track:", uri);
        // Your logic to play the track, such as updating the iframe URL or interacting with the Spotify API.
      },
    },
  };
  </script>
  
  <style scoped>
  /* Optional: Add some spacing and ensure alignment */
  .home-page {
    padding: 20px;
  }
  
  .grid {
    display: flex;
    justify-content: space-between;
  }
  
  .row {
    display: flex;
    flex-wrap: wrap;
  }
  
  .col-6 {
    width: 48%; /* Ensures each section takes up half the width */
    margin-bottom: 20px;
  }
  
  .col-4 {
    width: 30%; /* Each event and music card will take up 30% of the space */
    margin-bottom: 20px;
  }
  
  /* Additional styling if needed for spacing or appearance */
  </style>
  