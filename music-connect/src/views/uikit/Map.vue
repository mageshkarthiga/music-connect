<template>
  <div class="map-wrapper">
    <Button label="Load Map" icon="pi pi-map" @click="initMap" class="mb-3" />
    <div ref="map" class="map-container" v-show="mapLoaded"></div>

    <ul>
      <li v-for="(loc, userId) in userLocations" :key="userId">
        {{ userId }}: {{ loc.lat }}, {{ loc.lon }}
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { Loader } from "@googlemaps/js-api-loader";
import { getAllUserLocations } from "@/firebase/locationController";

const map = ref(null);
const mapLoaded = ref(false);
const userLocations = ref({});
let mapInstance;

// Google Maps loader
const loader = new Loader({
  apiKey: process.env.GOOGLE_API_KEY, // Replace with your actual API key
  version: "weekly",
  libraries: ["maps", "marker"],
});

// Static landmarks
const mockLandmarks = [
  { name: "SMU", lat: 1.28944, lng: 103.849983 },
  { name: "Springleaf Prata", lat: 1.401419, lng: 103.824061 },
  { name: "Cafe", lat: 1.355, lng: 103.812 },
];

const person_icon = document.createElement("img");
person_icon.src = "/demo/images/person_light.svg"; // or any valid path or URL
person_icon.style.width = "40px";
person_icon.style.height = "40px";

const other_people_icon = document.createElement("img");
other_people_icon.src = "/demo/images/person_dark.svg"; // or any valid path or URL
other_people_icon.style.width = "40px";
other_people_icon.style.height = "40px";

onMounted(async () => {
  userLocations.value = await getAllUserLocations();
  console.log("User Locations:", userLocations.value);
});

async function initMap() {
  try {
    await loader.load();

    if (!navigator.geolocation) {
      alert("Geolocation not supported");
      return;
    }

    navigator.geolocation.getCurrentPosition(
      async (pos) => {
        const userLocation = {
          lat: pos.coords.latitude,
          lng: pos.coords.longitude,
        };

        mapInstance = new google.maps.Map(map.value, {
          center: userLocation,
          zoom: 18,
          mapTypeId: "roadmap",
          mapId: "DEMO_MAP_ID",
        });

        // You
        new google.maps.marker.AdvancedMarkerElement({
          map: mapInstance,
          position: userLocation,
          title: "You are here",
          content: person_icon,
        });

        // Landmarks
        mockLandmarks.forEach((loc) => {
          new google.maps.marker.AdvancedMarkerElement({
            map: mapInstance,
            position: { lat: loc.lat, lng: loc.lng },
            title: `Landmark: ${loc.name}`,
          });
        });

        // Other users from Firestore
        Object.entries(userLocations.value).forEach(([userId, loc]) => {
          new google.maps.marker.AdvancedMarkerElement({
            map: mapInstance,
            position: {
              lat: loc.lat,
              lng: loc.lon,
            },
            title: `User: ${userId}`,
            content: other_people_icon,
          });
        });

        mapLoaded.value = true;
      },
      () => {
        alert("Unable to retrieve your location");
      }
    );
  } catch (error) {
    console.error("Google Maps failed to load:", error);
  }
}
</script>

<style scoped>
.map-wrapper {
  width: 100%;
  height: 100%;
  padding: 1rem;
}

.map-container {
  width: 100%;
  height: 80vh;
  border: 1px solid #ccc;
  border-radius: 8px;
}
</style>
