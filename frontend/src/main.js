import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import { createPinia } from "pinia";

import Aura from "@primevue/themes/aura";
import PrimeVue from "primevue/config";

import ConfirmationService from "primevue/confirmationservice";
import ToastService from "primevue/toastservice";

import { initLayoutFromFirestore } from "@/layout/composables/layoutController";
import { watchLayoutChanges } from "@/layout/composables/layoutController";
import "@/assets/styles.scss";

import { auth } from "@/firebase/firebase";
import { onAuthStateChanged } from "firebase/auth";
import {
  initLayoutFromFirestore,
  watchLayoutChanges,
} from "@/firebase/layoutController";
import {
  startLocationWatcher,
  stopLocationWatcher,
} from "@/firebase/locationController";
import { ref } from "vue";

const userLocation = ref({ lat: null, lon: null });

const app = createApp(App);
const pinia = createPinia();
app.use(pinia);
app.use(router);

// Use PrimeVue with configurations only once
app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: {
      darkModeSelector: ".app-dark",
    },
  },
});

await initLayoutFromFirestore();
watchLayoutChanges(); // Call the function to start watching

app.use(ToastService);
app.use(ConfirmationService);

onAuthStateChanged(auth, async (user) => {
  if (user) {
    await initLayoutFromFirestore(user.uid);
    watchLayoutChanges(user.uid);
    startLocationWatcher(user.uid, userLocation);
  } else {
    stopLocationWatcher();
    router.replace("/auth/login");
  }

  if (!app._instance) {
    app.mount("#app");
  }
});
