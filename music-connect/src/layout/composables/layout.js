import { computed, reactive, watch } from "vue";
import { readDocument, updateDocument } from "@/firebase/fireStoreService";
import { updateColors, applyPreset } from "@/layout/AppConfig";
import { primaryColors, surfaces } from "@/layout/AppConfig";
export const layoutConfig = reactive({
  preset: "Aura",
  primary: "emerald",
  surface: null,
  darkTheme: false,
  menuMode: "static",
});

const layoutState = reactive({
  staticMenuDesktopInactive: false,
  overlayMenuActive: false,
  profileSidebarVisible: false,
  configSidebarVisible: false,
  staticMenuMobileActive: false,
  menuHoverActive: false,
  activeMenuItem: null,
});

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

export function useLayout() {
  const setActiveMenuItem = (item) => {
    layoutState.activeMenuItem = item.value || item;
  };

  const toggleDarkMode = () => {
    if (!document.startViewTransition) {
      executeDarkModeToggle();
      return;
    }
    document.startViewTransition(() => executeDarkModeToggle(event));
  };

  const executeDarkModeToggle = () => {
    layoutConfig.darkTheme = !layoutConfig.darkTheme;
    document.documentElement.classList.toggle("app-dark");
  };

  const toggleMenu = () => {
    if (layoutConfig.menuMode === "overlay") {
      layoutState.overlayMenuActive = !layoutState.overlayMenuActive;
    }
    if (window.innerWidth > 991) {
      layoutState.staticMenuDesktopInactive =
        !layoutState.staticMenuDesktopInactive;
    } else {
      layoutState.staticMenuMobileActive = !layoutState.staticMenuMobileActive;
    }
  };

  const isSidebarActive = computed(
    () => layoutState.overlayMenuActive || layoutState.staticMenuMobileActive
  );

  const isDarkTheme = computed(() => layoutConfig.darkTheme);

  const getPrimary = computed(() => layoutConfig.primary);

  const getSurface = computed(() => layoutConfig.surface);

  return {
    layoutConfig,
    layoutState,
    toggleMenu,
    isSidebarActive,
    isDarkTheme,
    getPrimary,
    getSurface,
    setActiveMenuItem,
    toggleDarkMode,
  };
}
