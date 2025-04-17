<template>
  <div class="map-wrapper">
    <div class="flex justify-between items-center mb-3">
      <Button label="Load Map" icon="pi pi-map" @click="initMap" />
      <RouterLink to="/pages/search" class="text-blue-500 hover:underline">
        Go to Search
      </RouterLink>
    </div>

    <div ref="map" class="map-container" v-show="mapLoaded"></div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, onBeforeUnmount } from "vue";
import { getAllUserLocations } from "@/firebase/locationController";
import { API_BASE_URL } from "@/service/apiConfig";

const map = ref(null);
const mapLoaded = ref(false);
const userLocations = ref({});
let mapInstance = null;
let hasMapInitialised = false;
let currentPopup = null;

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

function createMarkerImage(src) {
  const wrapper = document.createElement("div");
  const img = document.createElement("img");
  img.src = src;
  img.style.width = "40px";
  img.style.height = "40px";
  img.style.cursor = "pointer";
  wrapper.appendChild(img);
  return wrapper;
}

function createPopupOverlay(lat, lng, title, description, image) {
  if (currentPopup) {
    currentPopup.setMap(null);
    currentPopup = null;
  }

  class PopupOverlay extends google.maps.OverlayView {
    div;

    constructor(position) {
      super();
      this.position = position;
    }

    onAdd() {
      this.div = document.createElement("div");
      this.div.style.position = "absolute";
      this.div.style.zIndex = 9999;
      this.div.innerHTML = `
        <div style="
          background: white;
          border-radius: 8px;
          box-shadow: 0 4px 8px rgba(0,0,0,0.3);
          padding: 12px;
          width: 250px;
          transition: all 0.3s ease;
        ">
          <img src="${image}" style="width:100%; height:120px; object-fit:cover; border-radius:4px; margin-bottom:8px;" />
          <h3 style="margin: 0 0 4px; font-size: 16px;">${title}</h3>
          <p style="margin: 0; font-size: 14px; color: #555;">${description}</p>
        </div>
      `;

      const panes = this.getPanes();
      panes.floatPane.appendChild(this.div);
    }

    draw() {
      const projection = this.getProjection();
      const point = projection.fromLatLngToDivPixel(this.position);
      if (point && this.div) {
        this.div.style.left = `${point.x - 125}px`;
        this.div.style.top = `${point.y - 180}px`;
      }
    }

    onRemove() {
      if (this.div && this.div.parentNode) {
        this.div.parentNode.removeChild(this.div);
      }
      this.div = null;
    }
  }

  currentPopup = new PopupOverlay(new google.maps.LatLng(lat, lng));
  currentPopup.setMap(mapInstance);
}

onMounted(async () => {
  userLocations.value = await getAllUserLocations();
  document.addEventListener("visibilitychange", handleVisibility);
  window.addEventListener("focus", handleVisibility);
});

onBeforeUnmount(() => {
  document.removeEventListener("visibilitychange", handleVisibility);
  window.removeEventListener("focus", handleVisibility);
});

function handleVisibility() {
  if (document.visibilityState === "visible") {
    ensureMapRedraw();
  }
}

function ensureMapRedraw() {
  if (mapInstance && map.value) {
    map.value.innerHTML = "";
    mapInstance = null;
    hasMapInitialised = false;
    currentPopup = null;
    initMapCallback();
  }
}

function loadGoogleMapsScript(callback) {
  if (window.google?.maps) return;

  const existingScript = document.querySelector(
    `script[src*="maps.googleapis.com/maps/api/js"]`
  );
  if (existingScript) return;

  const script = document.createElement("script");
  script.src = `${API_BASE_URL}/maps?libraries=maps,marker&callback=${callback}&loading=async`;
  script.async = true;
  script.defer = true;
  document.head.appendChild(script);
}

function initMapCallback() {
  if (hasMapInitialised || !navigator.geolocation || !map.value) return;

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
        mapId: "DEMO_MAP_ID",
      });

      mapInstance.addListener("click", () => {
        if (currentPopup) currentPopup.setMap(null);
        currentPopup = null;
      });

      // Self
      const youImg = createMarkerImage("/demo/images/person_light.svg");
      const selfMarker = new google.maps.marker.AdvancedMarkerElement({
        map: mapInstance,
        position: userLocation,
        title: "You are here",
        content: youImg,
      });
      selfMarker.addListener("click", () => {
        createPopupOverlay(
          userLocation.lat,
          userLocation.lng,
          "You",
          "This is your current location.",
          "/demo/images/person_light.svg"
        );
      });

      // Landmarks
      mockLandmarks.forEach((loc) => {
        const landmarkImg = createMarkerImage("/demo/images/logo.svg");
        const marker = new google.maps.marker.AdvancedMarkerElement({
          map: mapInstance,
          position: { lat: loc.lat, lng: loc.lng },
          title: `Landmark: ${loc.name}`,
          content: landmarkImg,
        });
        marker.addListener("click", () => {
          createPopupOverlay(
            loc.lat,
            loc.lng,
            loc.name,
            loc.description,
            loc.image
          );
        });
      });

      // Other users
      Object.entries(userLocations.value).forEach(([userId, loc]) => {
        const userImg = createMarkerImage("/demo/images/person_dark.svg");
        const marker = new google.maps.marker.AdvancedMarkerElement({
          map: mapInstance,
          position: { lat: loc.lat, lng: loc.lon },
          title: `User: ${userId}`,
          content: userImg,
        });
        marker.addListener("click", () => {
          createPopupOverlay(
            loc.lat,
            loc.lon,
            `User: ${userId}`,
            "Recently active in this area.",
            "/demo/images/person_dark.svg"
          );
        });
      });

      mapLoaded.value = true;
      hasMapInitialised = true;
    },
    () => {}
  );
}

async function initMap() {
  mapLoaded.value = false;
  await nextTick();
  map.value.innerHTML = "";
  mapInstance = null;
  hasMapInitialised = false;

  if (window.google?.maps) {
    initMapCallback();
  } else {
    window.initMapCallback = initMapCallback;
    loadGoogleMapsScript("initMapCallback");
  }
}
</script>

<style scoped>
.map-wrapper {
  width: 100%;
  height: 100%;
  padding: 1rem;
  position: relative;
}

.map-container {
  width: 100%;
  height: 80vh;
  border: 1px solid #ccc;
  border-radius: 8px;
}
</style>
