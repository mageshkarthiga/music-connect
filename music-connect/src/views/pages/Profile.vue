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
    style="max-width: 1000px; margin: 2rem auto; text-align: center"
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

    <Fluid>
      <div class="flex gap-4 mb-4 md:flex">
        <div class="md:w-1/2">
          <div class="card flex flex-col w-full">
            <div class="font-semibold text-xl mb-4">Skeleton</div>
            <div class="rounded-border border border-surface p-6">
              <div class="flex mb-4">
                <Skeleton shape="circle" size="4rem" class="mr-2"></Skeleton>
                <div>
                  <Skeleton width="10rem" class="mb-2"></Skeleton>
                  <Skeleton width="5rem" class="mb-2"></Skeleton>
                  <Skeleton height=".5rem"></Skeleton>
                </div>
              </div>
              <Skeleton width="100%" height="150px"></Skeleton>
              <div class="flex justify-between mt-4">
                <Skeleton width="4rem" height="2rem"></Skeleton>
                <Skeleton width="4rem" height="2rem"></Skeleton>
              </div>
            </div>
          </div>
        </div>
        <div class="md:w-1/2">
          <div class="card flex flex-col w-full">
            <div class="font-semibold text-xl mb-4">Skeleton</div>
            <div class="rounded-border border border-surface p-6">
              <div class="flex mb-4">
                <Skeleton shape="circle" size="4rem" class="mr-2"></Skeleton>
                <div>
                  <Skeleton width="10rem" class="mb-2"></Skeleton>
                  <Skeleton width="5rem" class="mb-2"></Skeleton>
                  <Skeleton height=".5rem"></Skeleton>
                </div>
              </div>
              <Skeleton width="100%" height="150px"></Skeleton>
              <div class="flex justify-between mt-4">
                <Skeleton width="4rem" height="2rem"></Skeleton>
                <Skeleton width="4rem" height="2rem"></Skeleton>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="flex gap-4 md:flex">
        <div class="md:w-1/2">
          <div class="card flex flex-col w-full">
            <div class="font-semibold text-xl mb-4">Skeleton</div>
            <div class="rounded-border border border-surface p-6">
              <div class="flex mb-4">
                <Skeleton shape="circle" size="4rem" class="mr-2"></Skeleton>
                <div>
                  <Skeleton width="10rem" class="mb-2"></Skeleton>
                  <Skeleton width="5rem" class="mb-2"></Skeleton>
                  <Skeleton height=".5rem"></Skeleton>
                </div>
              </div>
              <Skeleton width="100%" height="150px"></Skeleton>
              <div class="flex justify-between mt-4">
                <Skeleton width="4rem" height="2rem"></Skeleton>
                <Skeleton width="4rem" height="2rem"></Skeleton>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Fluid>
  </div>
</template>
