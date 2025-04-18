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
import { onAuthStateChanged, getAuth } from "firebase/auth";
import { getAllUserLocations } from "@/firebase/locationController";
import UserService from "@/service/UserService";
import EventService from "@/service/EventService";
import { API_BASE_URL } from "@/service/apiConfig";

const map = ref(null);
const mapLoaded = ref(false);
const userLocations = ref([]);
const eventLandmarks = ref([]);
let mapInstance = null;
let hasMapInitialised = false;
let currentPopup = null;

function createMarkerImage(src, isSelf = false, isEvent = false) {
  const wrapper = document.createElement("div");
  const img = document.createElement("img");
  img.src = src;
  img.style.width = "40px";
  img.style.height = "40px";
  img.style.cursor = "pointer";

  if (isEvent) {
    img.style.borderRadius = "0"; // square
    img.style.border = "1px solid black";
  } else {
    img.style.borderRadius = "50%"; // circle for users
    img.style.border = isSelf ? "3px solid #16a34a" : "1px solid black";
  }

  wrapper.appendChild(img);
  return wrapper;
}

function createPopupOverlay(lat, lng, title, description, image, link = null) {
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
      const buttonHTML = link
        ? `<a href="${link}" style="
              display:inline-block;
              margin-top:8px;
              background:#2563eb;
              color:white;
              padding:6px 12px;
              border-radius:4px;
              text-align:center;
              text-decoration:none;
              font-size:14px;
            ">View Profile</a>`
        : "";
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
          ${buttonHTML}
        </div>
      `;
      this.getPanes().floatPane.appendChild(this.div);
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

onMounted(() => {
  const auth = getAuth();
  onAuthStateChanged(auth, async (firebaseUser) => {
    if (!firebaseUser) return;

    const currentUID = firebaseUser.uid;

    const [userLocs, users, venues] = await Promise.all([
      getAllUserLocations(),
      UserService.getAllUsers(),
      EventService.getAllEventVenues(),
    ]);
    console.log("Venue data:", venues);

    const locatedUsers = users.filter((u) => userLocs[u.firebase_uid]);

    userLocations.value = locatedUsers.map((user) => {
      const { lat, lon } = userLocs[user.firebase_uid];
      return {
        id: user.user_id,
        name: user.user_name || `User: ${user.user_id}`,
        lat,
        lon,
        firebase_uid: user.firebase_uid,
        profilePhotoUrl: user.profile_photo_url,
        isSelf: user.firebase_uid === currentUID,
      };
    });

    const venueSet = new Map();
    venues.forEach((evt) => {
      evt.venues?.forEach((venue) => {
        if (venue.lat && venue.lon && !venueSet.has(venue.venue_id)) {
          venueSet.set(venue.venue_id, {
            lat: venue.lat,
            lon: venue.lon,
            name: venue.venue_name,
            description: venue.location,
            image:
              evt.event_image_url || venue.seat_map || "/demo/images/logo.svg",
          });
        }
      });
    });

    eventLandmarks.value = Array.from(venueSet.values());

    initMap();
  });

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

      userLocations.value.forEach((user) => {
        const imageSrc = user.profilePhotoUrl?.trim()
          ? user.profilePhotoUrl
          : "/demo/images/person_dark.svg";

        const userImg = createMarkerImage(imageSrc, user.isSelf);

        const marker = new google.maps.marker.AdvancedMarkerElement({
          map: mapInstance,
          position: { lat: user.lat, lng: user.lon },
          title: user.name,
          content: userImg,
        });

        marker.addListener("click", () => {
          createPopupOverlay(
            user.lat,
            user.lon,
            user.name,
            "Recently active in this area.",
            imageSrc,
            `/profile?user_id=${user.id}`
          );
        });
      });

      eventLandmarks.value.forEach((venue) => {
        const landmarkImg = createMarkerImage(venue.image, false, true);

        const marker = new google.maps.marker.AdvancedMarkerElement({
          map: mapInstance,
          position: { lat: venue.lat, lng: venue.lon },
          title: `Venue: ${venue.name}`,
          content: landmarkImg,
        });

        marker.addListener("click", () => {
          createPopupOverlay(
            venue.lat,
            venue.lon,
            venue.name,
            venue.description,
            venue.image
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
  mapLoaded.value = true;
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
