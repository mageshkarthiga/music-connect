// layoutController.js

import { watch } from "vue";
import { readDocument, updateDocument } from "@/firebase/fireStoreService";
import {
  primaryColors,
  surfaces,
  layoutConfig,
  layoutState,
  updateColors,
  applyPreset,
} from "@/layout/composables/stateConfig";

// Your existing logic...

export async function initLayoutFromFirestore() {
  try {
    const docData = await readDocument("layoutConfigs", "defaultLayout");
    if (docData) {
      // Existing logic to handle layoutConfig and layoutState
    }
  } catch (error) {
    console.error("Failed to load layout from Firestore:", error);
  }
}

// Example of a function you might want to export for watching layout changes
export function watchLayoutChanges() {
  watch(
    () => layoutConfig,
    (newVal) => {
      // Logic to watch changes and handle them
    },
    { deep: true }
  );

  watch(
    () => layoutState,
    (newVal) => {
      // Logic to watch changes and handle them
    },
    { deep: true }
  );
}
