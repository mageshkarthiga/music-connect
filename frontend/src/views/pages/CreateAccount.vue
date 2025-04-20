<script setup>
import { ref, watch, computed } from "vue";
import { useRouter } from "vue-router";
import PlaceAutoComplete from "@/components/map/PlaceAutoComplete.vue";
import FileUpload from "primevue/fileupload";
import InputGroup from "primevue/inputgroup";
import InputGroupAddon from "primevue/inputgroupaddon";
import InputText from "primevue/inputtext";
import Button from "primevue/button";
import FloatingConfigurator from "@/components/FloatingConfigurator.vue";
import { useLayout } from "@/layout/composables/stateConfig";
import userService from "@/service/UserService";
import {
  ref as storageRef,
  uploadBytes,
  getDownloadURL,
} from "firebase/storage";
import { storage } from "@/firebase/firebase";
import { getAuth, onAuthStateChanged } from "firebase/auth";

const auth = getAuth();

const { isDarkTheme } = useLayout();
const router = useRouter();

const defaultProfilePhotoUrl = computed(() =>
  isDarkTheme.value
    ? "/demo/images/person_dark.svg"
    : "/demo/images/person_light.svg"
);
const username = ref("");
const phoneNumber = ref("");
const profilePhoto = ref(null);
const profilePhotoUrl = ref(defaultProfilePhotoUrl.value);
const selectedLocation = ref("");
const authError = ref("");

var email = "";
var fb_id = "";
onAuthStateChanged(auth, async (user) => {
  if (!user) {
    router.push("/auth/login");
    return;
  }

  email = user.email;
  fb_id = user.uid;
  localStorage.setItem("uid", fb_id);

  if (!email) {
    console.error("Email not found in query params.");
    router.push("/auth/login");
  }
  if (!fb_id) {
    console.error("firebase uid not found in query params.");
    router.push("/auth/login");
  }
});

const handleFileUpload = (event) => {
  if (event.files && event.files.length > 0) {
    profilePhoto.value = event.files[0];
    profilePhotoUrl.value = URL.createObjectURL(event.files[0]);
  }
};

watch(isDarkTheme, () => {
  if (!profilePhoto.value) {
    profilePhotoUrl.value = defaultProfilePhotoUrl.value;
  }
});

const handlePlaceSelected = (place) => {
  selectedLocation.value = place;
};

const uploadProfilePhoto = async (file) => {
  const auth = getAuth();
  const user = auth.currentUser;

  if (!user) {
    console.error("User not found.");
    router.push("/auth/login");
  }
  const uid = user.uid;
  const filePath = `profile_photos/${uid}/${file.name}`;
  const fileRef = storageRef(storage, filePath);

  await uploadBytes(fileRef, file);
  return await getDownloadURL(fileRef);
};

const handleSubmit = async () => {
  const errors = [];

  if (!username.value) errors.push("Username is missing.");
  if (!phoneNumber.value) errors.push("Phone number is missing.");
  if (phoneNumber.value && !+phoneNumber.value) errors.push("Phone number is not valid.");
  if (!selectedLocation.value) errors.push("Location is not selected.");

  if (errors.length) {
    authError.value = errors.join(" ");
    return;
  }

  let photoURL = "";

  try {
    if (profilePhoto.value) {
      photoURL = await uploadProfilePhoto(profilePhoto.value);
    }

    const response = await userService.createUser({
      userName: username.value,
      phoneNumber: phoneNumber.value,
      location: selectedLocation.value,
      emailAddress: email,
      profilePhotoUrl: photoURL,
      firebaseUID: fb_id,
    });
    localStorage.setItem("uid", response.userID);
  } catch (error) {
    ("CREATE USER FAILED");
  }

  router.push("/pages/home");
};
</script>

<template>
  <FloatingConfigurator />
  <div
    class="bg-surface-50 dark:bg-surface-950 flex items-center justify-center min-h-screen min-w-[100vw] overflow-hidden"
  >
    <div class="flex flex-col items-center w-full">
      <div class="card p-6 w-full max-w-lg">
        <h2 class="text-2xl font-semibold mb-4 text-center">
          Set Up Your Account
        </h2>
        <img
          :src="profilePhotoUrl"
          alt="Profile Photo"
          class="w-32 h-32 rounded-full object-cover mx-auto mb-4"
        />
        <div class="flex flex-col gap-4">
          <FileUpload
            mode="basic"
            accept="image/*"
            :maxFileSize="1000000"
            @select="handleFileUpload"
            chooseLabel="Upload Profile Photo"
          />
          <InputGroup>
            <InputGroupAddon>
              <i class="pi pi-user"></i>
            </InputGroupAddon>
            <InputText v-model="username" placeholder="Username" />
          </InputGroup>
          <InputGroup>
            <InputGroupAddon>
              <i class="pi pi-phone"></i>
            </InputGroupAddon>
            <InputText
              v-model="phoneNumber"
              placeholder="Phone Number"
              type="tel"
            />
          </InputGroup>
          <InputGroup>
            <InputGroupAddon>
              <i class="pi pi-map-marker"></i>
            </InputGroupAddon>
            <PlaceAutoComplete
              @place-selected="handlePlaceSelected"
              v-model="selectedLocation"
            />
          </InputGroup>
          <Message v-if="authError" severity="error" class="mb-2">
            {{ authError }}
          </Message>
          <Button
            label="Continue"
            @click="handleSubmit"
            class="p-button-success w-full"
          />
        </div>
      </div>
    </div>
  </div>
</template>
