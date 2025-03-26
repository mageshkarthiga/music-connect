<template>
  <div class="map-wrapper">
    <Button label="Load Map" icon="pi pi-map" @click="initMap" class="mb-3" />
    <div ref="map" class="map-container" v-show="mapLoaded"></div>
  </div>
</template>

<script setup>
import { ref } from "vue";

const map = ref(null);
const mapLoaded = ref(false);
let mapInstance;

const mockLandmarks = [
  { name: "SMU", lat: 1.28944, lng: 103.849983 },
  { name: "Springleaf Prata", lat: 1.401419, lng: 103.824061 },
  { name: "Cafe", lat: 1.355, lng: 103.812 },
];

const mockUsers = [
  { name: "SMU", lat: 1.28935, lng: 103.849984 },
  { name: "Springleaf Prata", lat: 1.401108, lng: 103.824061 },
  { name: "Cafe", lat: 1.315, lng: 103.812 },
];

function loadGoogleMapsScript() {
  return new Promise((resolve, reject) => {
    if (window.google && window.google.maps) return resolve();
    const script = document.createElement("script");
    script.src = "http://localhost:3000/api/maps";
    script.async = true;
    script.defer = true;
    script.onload = resolve;
    script.onerror = reject;
    document.head.appendChild(script);
  });
}

async function initMap() {
  try {
    await loadGoogleMapsScript();

    if (!navigator.geolocation) {
      alert("Geolocation not supported");
      return;
    }

    navigator.geolocation.getCurrentPosition(
      (pos) => {
        const userLocation = {
          lat: pos.coords.latitude,
          lng: pos.coords.longitude,
        };

        mapInstance = new google.maps.Map(map.value, {
          center: userLocation,
          zoom: 18,
          mapTypeId: "roadmap",
        });

        // User avatar
        new google.maps.Marker({
          position: userLocation,
          map: mapInstance,
          icon: {
            url: "/demo/images/person_light.svg",
            scaledSize: new google.maps.Size(40, 40),
          },
          title: "You are here",
        });

        // Mock landmarks
        mockLandmarks.forEach((loc) => {
          new google.maps.Marker({
            position: { lat: loc.lat, lng: loc.lng },
            map: mapInstance,
            title: loc.name,
            icon: {
              url: "/demo/images/logo.svg",
              scaledSize: new google.maps.Size(40, 40),
            },
          });
        });

        // Mock users
        mockUsers.forEach((loc) => {
          new google.maps.Marker({
            position: { lat: loc.lat, lng: loc.lng },
            map: mapInstance,
            title: loc.name,
            icon: {
              url: "/demo/images/person_light.svg",
              scaledSize: new google.maps.Size(40, 40),
            },
          });
        });

        mapLoaded.value = true;
      },
      () => {
        alert("Unable to retrieve your location");
      }
    );
  } catch (e) {
    console.error("Google Maps failed to load", e);
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
