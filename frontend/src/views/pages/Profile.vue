<script>
import { defineComponent, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import UserService from "@/service/UserService";
import { getAuth, onAuthStateChanged } from "firebase/auth";

export default defineComponent({
  name: "Profile",
  setup() {
    const route = useRoute();
    const user = ref(null);
    const loading = ref(true);

    onMounted(async () => {
      const userId = Number(route.query.user_id);
      if (!isNaN(userId)) {
        // Public profile of others
        try {
          user.value = await UserService.getUser(userId);
          console.log(user.value);
        } catch (error) {
          console.error("Error fetching user:", error);
        }
      } else {
        // Personal profile
        const auth = getAuth();
        onAuthStateChanged(auth, async (firebaseUser) => {
          if (firebaseUser) {
            try {
              user.value = await UserService.getUserByFirebaseUID(
                firebaseUser.uid
              );
            } catch (error) {
              console.error("Error fetching user by Firebase UID:", error);
            }
          } else {
            console.warn("No Firebase user signed in.");
          }
        });
      }
      loading.value = false;
    });

    return {
      user,
      loading,
    };
  },
});
</script>

<template>
  <!-- Inline styling can be awkward, but it's requested here. -->
  <div
    class="profile-page"
    style="max-width: 500px; margin: 2rem auto; text-align: center"
  >
    <div
      v-if="loading"
      class="loading"
      style="font-size: 1.2rem; padding: 2rem"
    >
      <p>Loading...</p>
    </div>
    <div v-else-if="user" class="profile-details p-card p-p-4 p-shadow-4">
      <img
        :src="user.ProfilePhotoUrl"
        alt="Profile Photo"
        class="profile-photo"
        style="
          width: 120px;
          height: 120px;
          object-fit: cover;
          border-radius: 50%;
          border: 3px solid var(--primary-color);
        "
      />
      <h1 class="p-mt-3">{{ user.userName }}</h1>
      <p><strong>Email:</strong> {{ user.EmailAddress }}</p>
      <p><strong>Phone:</strong> {{ user.PhoneNumber }}</p>
      <p><strong>Location:</strong> {{ user.Location }}</p>
    </div>
    <div v-else class="error p-error" style="font-size: 1.2rem; padding: 2rem">
      <p>User not found.</p>
    </div>
  </div>
</template>
