import { updateDocument } from "@/firebase/fireStoreService";
import { collection, getDocs } from "firebase/firestore";
import { db } from "@/firebase/firebase";
let locationWatcherId = null;

export function startLocationWatcher(userId, locationRef) {
  if (!navigator.geolocation) {
    console.warn("Geolocation is not supported by this browser.");
    return;
  }

  locationWatcherId = navigator.geolocation.watchPosition(
    async (position) => {
      const coords = {
        lat: position.coords.latitude,
        lon: position.coords.longitude,
      };

      if (
        locationRef.value.lat !== coords.lat ||
        locationRef.value.lon !== coords.lon
      ) {
        locationRef.value = coords;
        try {
          await updateDocument("location", userId, coords);
        } catch (error) {
          console.error("Failed to update location in Firestore:", error);
        }
      }
    },
    (error) => {
      console.warn("Geolocation watch error:", error);
    },
    {
      enableHighAccuracy: true,
      maximumAge: 10000,
      timeout: 10000,
    }
  );
}

export function stopLocationWatcher() {
  if (locationWatcherId !== null) {
    navigator.geolocation.clearWatch(locationWatcherId);
    locationWatcherId = null;
  }
}

export async function getAllUserLocations() {
  try {
    const snapshot = await getDocs(collection(db, "location"));
    const locations = {};

    snapshot.forEach((doc) => {
      const data = doc.data();
      if (data.lat != null && data.lon != null) {
        locations[doc.id] = {
          lat: data.lat,
          lon: data.lon,
        };
      }
    });

    return locations;
  } catch (error) {
    console.error("Error fetching user locations:", error);
    return [];
  }
}
