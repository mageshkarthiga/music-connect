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

const COLLECTION_NAME = "layoutConfigs";
const DOC_ID = "defaultLayout";

export async function initLayoutFromFirestore() {
  try {
    const docData = await readDocument(COLLECTION_NAME, DOC_ID);
    if (docData) {
      if (docData.layoutConfig) {
        Object.assign(layoutConfig, docData.layoutConfig);
        if (layoutConfig.darkTheme) {
          document.documentElement.classList.add("app-dark");
        }
        const foundPrimary = primaryColors.value.find(
          (c) => c.name === layoutConfig.primary
        );
        if (foundPrimary) updateColors("primary", foundPrimary);

        const foundSurface = surfaces.value.find(
          (s) => s.name === layoutConfig.surface
        );
        if (foundSurface) updateColors("surface", foundSurface);
        if (layoutConfig.preset) applyPreset();
      }
      if (docData.layoutState) {
        Object.assign(layoutState, docData.layoutState);
      }
    }
  } catch (error) {
    console.error("Failed to load layout from Firestore:", error);
  }
}

// Watchers to save changes to Firestore
watch(
  () => layoutConfig,
  async (newVal) => {
    try {
      await updateDocument(COLLECTION_NAME, DOC_ID, {
        layoutConfig: { ...newVal },
      });
    } catch (error) {
      console.error("Failed to update layoutConfig in Firestore:", error);
    }
  },
  { deep: true }
);

watch(
  () => layoutState,
  async (newVal) => {
    try {
      await updateDocument(COLLECTION_NAME, DOC_ID, {
        layoutState: { ...newVal },
      });
    } catch (error) {
      console.error("Failed to update layoutState in Firestore:", error);
    }
  },
  { deep: true }
);
