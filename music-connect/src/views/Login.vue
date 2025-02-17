<!-- src/views/Login.vue -->
<template>
  <div class="p-d-flex p-jc-center p-ai-center" style="min-height: 100vh">
    <div class="p-shadow-3 p-p-4" style="width: 300px">
      <h2 class="p-text-center">Login</h2>
      <form @submit.prevent="login">
        <div class="p-field">
          <label for="email">Email:</label>
          <InputText id="email" v-model="email" type="email" required />
        </div>
        <div class="p-field">
          <label for="password">Password:</label>
          <Password id="password" v-model="password" toggleMask required />
        </div>
        <Button label="Login" type="submit" class="p-mt-2" />
      </form>
      <p class="p-text-center p-mt-2">
        Don't have an account?
        <router-link to="/signup">Sign up here</router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { auth } from "../firebase/firebase";
import { signInWithEmailAndPassword } from "firebase/auth";
import { useRouter } from "vue-router";

// PrimeVue components
import InputText from "primevue/inputtext";
import Password from "primevue/password";
import Button from "primevue/button";

const router = useRouter();
const email = ref("");
const password = ref("");

const login = async () => {
  try {
    await signInWithEmailAndPassword(auth, email.value, password.value);
    router.push("/"); // Redirect to Home on successful login
  } catch (error) {
    console.error("Login error:", error);
  }
};
</script>
