<!-- src/views/Signup.vue -->
<template>
  <div class="p-d-flex p-jc-center p-ai-center" style="min-height: 100vh">
    <div class="p-shadow-3 p-p-4" style="width: 300px">
      <h2 class="p-text-center">Sign Up</h2>
      <form @submit.prevent="signup">
        <div class="p-field">
          <label for="email">Email:</label>
          <InputText id="email" v-model="email" type="email" required />
        </div>
        <div class="p-field">
          <label for="password">Password:</label>
          <Password id="password" v-model="password" toggleMask required />
        </div>
        <Button label="Sign Up" type="submit" class="p-mt-2" />
      </form>
      <p class="p-text-center p-mt-2">
        Already have an account?
        <router-link to="/login">Login here</router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { auth } from "../firebase/firebase";
import { createUserWithEmailAndPassword } from "firebase/auth";
import { useRouter } from "vue-router";

// PrimeVue components
import InputText from "primevue/inputtext";
import Password from "primevue/password";
import Button from "primevue/button";

const router = useRouter();
const email = ref("");
const password = ref("");

const signup = async () => {
  try {
    await createUserWithEmailAndPassword(auth, email.value, password.value);
    router.push("/"); // Redirect to Home on successful signup
  } catch (error) {
    console.error("Signup error:", error);
  }
};
</script>
