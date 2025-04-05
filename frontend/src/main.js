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
watchLayoutChanges();  // Call the function to start watching

app.use(ToastService);
app.use(ConfirmationService);

app.mount("#app");
