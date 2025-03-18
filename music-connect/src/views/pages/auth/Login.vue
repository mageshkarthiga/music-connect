<script setup>
import FloatingConfigurator from "@/components/FloatingConfigurator.vue";
import { ref } from "vue";
import { auth } from "@/firebase/firebase";
import {
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
} from "firebase/auth";
import { useRouter } from "vue-router";

const email = ref("");
const password = ref("");
const authError = ref("");
const mode = ref("login");
const router = useRouter();

const login = async () => {
  authError.value = "";
  try {
    await signInWithEmailAndPassword(auth, email.value, password.value);
    router.push("/");
  } catch (error) {
    console.error("Login error:", error);
    switch (error.code) {
      case "auth/user-not-found":
        authError.value = "No account found with this email.";
        break;
      case "auth/invalid-email":
        authError.value = "Please enter a valid email address.";
        break;
      case "auth/invalid-password":
      case "auth/wrong-password":
        authError.value = "Incorrect password. Please try again.";
        break;
      case "auth/too-many-requests":
        authError.value = "Too many failed attempts. Try again later.";
        break;
      case "auth/operation-not-allowed":
        authError.value =
          "Email/password login is disabled for this application.";
        break;
      case "auth/internal-error":
        authError.value =
          "An unexpected error occurred. Please try again later.";
        break;
      default:
        authError.value =
          "Login failed. Please check your credentials and try again.";
    }
  }
};

const signUp = async () => {
  authError.value = "";
  try {
    await createUserWithEmailAndPassword(auth, email.value, password.value);
    router.push("/createaccount");
  } catch (error) {
    console.error("Sign Up error:", error);
    switch (error.code) {
      case "auth/email-already-exists":
      case "auth/email-already-in-use":
        authError.value =
          "This email is already in use. Please try logging in instead.";
        break;
      case "auth/invalid-email":
        authError.value = "Please enter a valid email address.";
        break;
      case "auth/weak-password":
      case "auth/invalid-password":
        authError.value = "Your password must be at least 6 characters long.";
        break;
      case "auth/missing-email":
        authError.value = "Email is required to create an account.";
        break;
      case "auth/missing-password":
        authError.value = "Password is required to create an account.";
        break;
      case "auth/operation-not-allowed":
        authError.value =
          "Sign-up is currently disabled. Please contact support.";
        break;
      case "auth/internal-error":
        authError.value =
          "An unexpected error occurred. Please try again later.";
        break;
      case "auth/too-many-requests":
        authError.value = "Too many attempts. Please try again later.";
        break;
      default:
        authError.value =
          "Sign-up failed. Please check your details and try again.";
    }
  }
};

const onSubmit = async () => {
  if (mode.value === "login") {
    await login();
  } else {
    await signUp();
  }
};

const toggleMode = () => {
  authError.value = "";
  mode.value = mode.value === "login" ? "signup" : "login";
};
</script>

<template>
  <FloatingConfigurator />
  <div
    class="bg-surface-50 dark:bg-surface-950 flex items-center justify-center min-h-screen min-w-[100vw] overflow-hidden"
  >
    <div class="flex flex-col items-center justify-center">
      <div
        class="rounded-[56px] p-[0.3rem] bg-gradient-to-b from-[var(--primary-color)] via-transparent to-transparent w-full"
      >
        <div
          class="w-full bg-surface-0 dark:bg-surface-900 py-20 px-8 md:w-[30rem] lg:w-[35rem] sm:px-20 rounded-[53px]"
        >
          <div class="w-full flex flex-col items-center">
            <img
              src="/demo/images/logo.svg"
              alt="Logo"
              class="mb-8"
              width="30%"
            />
            <div
              class="text-surface-900 dark:text-surface-0 text-3xl font-medium mb-4"
              v-if="mode === 'login'"
            >
              Welcome back
            </div>
            <div
              class="text-surface-900 dark:text-surface-0 text-3xl font-medium mb-4"
              v-else
            >
              Let's get started
            </div>
            <span
              class="text-muted-color font-medium mb-2"
              v-if="mode === 'login'"
            >
              Sign in to continue
            </span>
            <span class="text-muted-color font-medium mb-2" v-else>
              Sign up to continue
            </span>
          </div>

          <div>
            <form @submit.prevent="onSubmit">
              <label
                for="email1"
                class="block text-surface-900 dark:text-surface-0 text-xl font-medium mb-2 ml-1"
              >
                Email
              </label>
              <InputText
                id="email1"
                type="text"
                placeholder="Email address"
                class="w-full mb-8"
                v-model="email"
              />

              <label
                for="password1"
                class="block text-surface-900 dark:text-surface-0 font-medium text-xl mb-2 ml-1"
              >
                Password
              </label>
              <Password
                id="password1"
                v-model="password"
                placeholder="Password"
                :toggleMask="true"
                class="mb-8"
                fluid
                :feedback="false"
              />

              <Message v-if="authError" severity="error" class="mb-8">
                {{ authError }}
              </Message>
              <Button
                :label="mode === 'login' ? 'Sign In' : 'Sign Up'"
                class="w-full mb-4"
                type="submit"
              ></Button>
            </form>

            <div class="flex items-center justify-center mb-8 gap-1">
              <span>{{
                mode === "login"
                  ? "Don't have an account?"
                  : "Already have an account?"
              }}</span>
              <Button
                variant="link"
                @click="toggleMode"
                unstyled
                class="font-medium cursor-pointer underline text-primary"
              >
                {{ mode === "login" ? "Sign up now" : "Sign in now" }}
              </Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
