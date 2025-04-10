import { watch } from "vue";
import { readDocument, updateDocument } from "@/firebase/fireStoreService";
import {
  layoutConfig,
  layoutState,
  primaryColors,
  surfaces,
  updateColors,
  applyPreset,
} from "@/layout/composables/stateConfig";
import { toRaw } from "vue";

const COLLECTION_NAME = "layoutConfigs";

export async function initLayoutFromFirestore(userId) {
  if (!userId) return;
  try {
    const docData = await readDocument(COLLECTION_NAME, userId);
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

export function watchLayoutChanges(userId) {
  watch(
    () => layoutConfig,
    async (newVal) => {
      if (!userId) return;
      try {
        await updateDocument(COLLECTION_NAME, userId, {
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
      const rawLayoutState = toRaw(newVal); // Get the raw object

      console.log("Raw layoutState:", rawLayoutState);

      try {
        if (rawLayoutState) {
          await updateDocument(COLLECTION_NAME, userId, {
            layoutState: { ...rawLayoutState },
          });
        }
      } catch (error) {
        console.error("Failed to update layoutState in Firestore:", error);
      }
    },
    { deep: true }
  );
}
