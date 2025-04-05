<template>
  <div class="map-wrapper">
    <Button label="Load Map" icon="pi pi-map" @click="initMap" class="mb-3" />
    <div ref="map" class="map-container" v-show="mapLoaded"></div>

    <Popover ref="overlay" :dismissable="true" style="width: 450px">
      <Card class="w-72">
        <template #header>
          <img
            :src="cardData.image"
            alt="Profile"
            class="w-full h-32 object-cover"
          />
        </template>
        <template #title>
          {{ cardData.title }}
        </template>
        <template #content>
          <p>{{ cardData.description }}</p>
        </template>
      </Card>
    </Popover>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { Loader } from "@googlemaps/js-api-loader";
import { getAllUserLocations } from "@/firebase/locationController";

const map = ref(null);
const mapLoaded = ref(false);
const userLocations = ref({});
const overlay = ref();
const cardData = ref({ title: "", description: "", image: "" });
let mapInstance;

const loader = new Loader({
  apiKey: process.env.GOOGLE_API_KEY,
  version: "weekly",
  libraries: ["maps", "marker"],
});

const mockLandmarks = [
  {
    name: "SMU",
    lat: 1.28944,
    lng: 103.849983,
    description: "Singapore Management University – heart of city campus life.",
    image: "/demo/images/smu.jpg",
  },
  {
    name: "Springleaf Prata",
    lat: 1.401419,
    lng: 103.824061,
    description: "Late-night favourite for crispy prata and curry.",
    image: "/demo/images/prata.jpg",
  },
  {
    name: "Tiong Bahru Bakery",
    lat: 1.2842,
    lng: 103.828,
    description: "Famed for its buttery croissants and rustic charm.",
    image: "/demo/images/bakery.jpg",
  },
  {
    name: "Marina Bay Sands",
    lat: 1.2834,
    lng: 103.8607,
    description: "Iconic hotel with infinity pool and rooftop views.",
    image: "/demo/images/mbs.jpg",
  },
  {
    name: "HortPark",
    lat: 1.2753,
    lng: 103.7996,
    description: "Nature-lovers' hub for walks and gardening.",
    image: "/demo/images/hortpark.jpg",
  },
  {
    name: "Cafe",
    lat: 1.355,
    lng: 103.812,
    description: "Hidden gem for quiet study and rich espresso.",
    image: "/demo/images/cafe.jpg",
  },
];

onMounted(async () => {
  userLocations.value = await getAllUserLocations();
});

function showPopover(event, title, description, image) {
  cardData.value = { title, description, image };
  overlay.value.show(event);
}

async function initMap() {
  try {
    await loader.load();

    if (!navigator.geolocation) return;

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
        function createMarkerImage(src) {
          const img = document.createElement("img");
          img.src = src;
          img.style.width = "40px";
          img.style.height = "40px";
          return img;
        }

        const youMarker = new google.maps.marker.AdvancedMarkerElement({
          map: mapInstance,
          position: userLocation,
          title: "You are here",
          content: createMarkerImage("/demo/images/person_light.svg"),
        });

        youMarker.addListener("click", (e) => {
          showPopover(
            e.domEvent,
            "You",
            "This is your current location.",
            "/demo/images/person_light.svg"
          );
        });

        mockLandmarks.forEach((loc) => {
          const marker = new google.maps.marker.AdvancedMarkerElement({
            map: mapInstance,
            position: { lat: loc.lat, lng: loc.lng },
            title: `Landmark: ${loc.name}`,
            content: createMarkerImage("/demo/images/logo.svg"),
          });

          marker.addListener("click", (e) => {
            showPopover(e.domEvent, loc.name, loc.description, loc.image);
          });
        });

        Object.entries(userLocations.value).forEach(([userId, loc]) => {
          const person = new google.maps.marker.AdvancedMarkerElement({
            map: mapInstance,
            position: {
              lat: loc.lat,
              lng: loc.lon,
            },
            title: `User: ${userId}`,
            content: createMarkerImage("/demo/images/person_dark.svg"),
          });

          person.addListener("click", (e) => {
            showPopover(
              e.domEvent,
              `User: ${userId}`,
              "Recently active in this area.",
              "/demo/images/person_dark.svg"
            );
          });
        });

        mapLoaded.value = true;
      },
      () => {}
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
