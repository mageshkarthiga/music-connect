import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";

import Aura from "@primevue/themes/aura";
import PrimeVue from "primevue/config";
import ConfirmationService from "primevue/confirmationservice";
import ToastService from "primevue/toastservice";
import "@/assets/styles.scss";

import { auth } from "@/firebase/firebase";
import { onAuthStateChanged } from "firebase/auth";
import {
  initLayoutFromFirestore,
  watchLayoutChanges,
} from "@/layout/composables/layoutController";

const app = createApp(App);

app.use(router);
app.use(PrimeVue, {
  theme: {
    preset: Aura,
    options: {
      darkModeSelector: ".app-dark",
    },
  },
});
app.use(ToastService);
app.use(ConfirmationService);

onAuthStateChanged(auth, async (user) => {
  if (user) {
    await initLayoutFromFirestore(user.uid);
    watchLayoutChanges(user.uid);
  } else {
    router.replace("/auth/login");
  }

  if (!app._instance) {
    app.mount("#app");
  }
});
