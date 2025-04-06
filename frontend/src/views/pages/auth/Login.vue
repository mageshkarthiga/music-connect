<script setup>
import FloatingConfigurator from "@/components/FloatingConfigurator.vue";
import { ref } from "vue";
import { auth } from "@/firebase/firebase";
import {
  signInWithEmailAndPassword,
  createUserWithEmailAndPassword,
  signInWithPopup,
  GoogleAuthProvider,
  getAdditionalUserInfo,
} from "firebase/auth";


import supabaseUserService from "@/service/supabaseUserService";
import { supabase } from "@/service/supabaseClient";


import { useRouter } from "vue-router";

const email = ref("");
const password = ref("");
const authError = ref("");
const mode = ref("login");
const router = useRouter();


async function storeUserInSupabase(firebaseUID) {
  try {
    const user = auth.currentUser;
    console.log("Firebase user:", user);

    if (!user) throw new Error('No user logged in with Firebase');

    const user_name = user.displayName || null;
    const email_address = user.email || null;
    const phone_number = user.phoneNumber || null;
    const profile_photo_url = user.photoURL || null;
    const location = null; // Set this if you have it from elsewhere


    console.log('Storing user in Supabase:', {
      firebaseUID,
      user_name,
      email_address,
      phone_number,
      location,
      profile_photo_url,
    });

    const { data, error } = await supabase
      .from('users')
      .insert([
        {
          firebase_uid: firebaseUID,
          user_name: user_name,
          email_address: email_address,
          phone_number: phone_number,
          location: location,
          profile_photo_url: profile_photo_url,
        }
      ]);

    if (error) {
      throw new Error(`Error inserting user into Supabase: ${error.message}`);
    }

    
    console.log('✅ User stored successfully in Supabase:');
  } catch (error) {
    console.error('❌ Error storing user in Supabase:', error);
    throw error;
  }
}


const login = async () => {
  authError.value = "";
  try {
    await signInWithEmailAndPassword(auth, email.value, password.value);
    const userCredential = await signInWithEmailAndPassword(auth, email.value, password.value);
    await storeUserInSupabase(userCredential.user.uid); 
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
      case "auth/invalid-credential":
      case "auth/invalid-password":
      case "auth/wrong-password":
        authError.value = "Incorrect email or password. Please try again.";
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
    const userCredential = await signInWithEmailAndPassword(auth, email.value, password.value);
    await storeUserInSupabase(userCredential.user.uid); 

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

const googleProvider = new GoogleAuthProvider();

const signInWithGoogle = async () => {
  authError.value = "";
  try {
    const userCredential = await signInWithPopup(auth, googleProvider);
    const additionalUserInfo = getAdditionalUserInfo(userCredential);
    await storeUserInSupabase(userCredential.user.uid); 

    if (additionalUserInfo?.isNewUser) {
      router.push("/createaccount");
    } else {
      router.push("/");
    }
  } catch (error) {
    console.error("Google Sign-In error:", error);
    if (error.code === "auth/popup-closed-by-user") {
      authError.value = "Google Sign-In cancelled.";
    } else if (error.code === "auth/cancelled-popup-request") {
      return;
    } else if (error.code === "auth/account-exists-with-different-credential") {
      authError.value =
        "An account already exists with this email using a different sign-in method. Try logging in with that method.";
    } else {
      authError.value = "Google Sign-In failed. Please try again.";
    }
  }
};
</script>

<template>
  <FloatingConfigurator />
  <div
    class="bg-surface-50 dark:bg-surface-950 flex items-center justify-center min-h-screen min-w-[100vw] overflow-hidden py-10"
  >
    <div class="flex flex-col items-center justify-center">
      <div
        class="rounded-[56px] p-[0.3rem] bg-gradient-to-b from-[var(--primary-color)] via-transparent to-transparent w-full"
      >
        <div
          class="w-full bg-surface-0 dark:bg-surface-900 py-12 px-8 md:w-[30rem] lg:w-[35rem] sm:px-16 rounded-[53px]"
        >
          <div class="w-full flex flex-col items-center mb-8">
            <img
              src="/demo/images/logo.svg"
              alt="Logo"
              class="mb-6"
              width="30%"
            />
            <div
              class="text-surface-900 dark:text-surface-0 text-3xl font-medium mb-2"
            >
              {{ mode === "login" ? "Welcome back" : "Let's get started" }}
            </div>
            <span class="text-muted-color font-medium">
              {{
                mode === "login" ? "Sign in to continue" : "Sign up to continue"
              }}
            </span>
          </div>

          <Message
            v-if="authError"
            severity="error"
            class="mb-6"
            :closable="false"
          >
            {{ authError }}
          </Message>

          <form @submit.prevent="onSubmit">
            <div class="flex flex-col gap-6">
              <div>
                <label
                  for="email1"
                  class="block text-surface-900 dark:text-surface-0 text-base font-medium mb-2 ml-1"
                >
                  Email
                </label>
                <InputText
                  id="email1"
                  type="text"
                  placeholder="Email address"
                  class="w-full"
                  v-model="email"
                  aria-describedby="email-help"
                  :invalid="
                    !!authError &&
                    (authError.includes('email') ||
                      authError.includes('credential') ||
                      authError.includes('account'))
                  "
                />
              </div>

              <div>
                <label
                  for="password1"
                  class="block text-surface-900 dark:text-surface-0 font-medium text-base mb-2 ml-1"
                >
                  Password
                </label>
                <Password
                  id="password1"
                  v-model="password"
                  placeholder="Password"
                  :toggleMask="true"
                  class="w-full"
                  fluid
                  :feedback="mode === 'signup'"
                  inputClass="w-full"
                  :invalid="
                    !!authError &&
                    (authError.includes('password') ||
                      authError.includes('credential'))
                  "
                />
              </div>

              <Button
                :label="mode === 'login' ? 'Sign In' : 'Sign Up'"
                class="w-full"
                type="submit"
                :loading="false"
              ></Button>
            </div>
          </form>

          <div class="flex items-center my-6">
            <div
              class="h-px bg-surface-200 dark:bg-surface-700 flex-grow"
            ></div>
            <span class="mx-4 text-muted-color text-sm font-medium">OR</span>
            <div
              class="h-px bg-surface-200 dark:bg-surface-700 flex-grow"
            ></div>
          </div>

          <Button
            label="Sign in with Google"
            icon="pi pi-google"
            severity="contrast"
            class="w-full mb-6"
            @click="signInWithGoogle"
          ></Button>

          <div class="flex items-center justify-center gap-1">
            <span class="text-muted-color">{{
              mode === "login"
                ? "Don't have an account?"
                : "Already have an account?"
            }}</span>
            <Button
              variant="link"
              @click="toggleMode"
              unstyled
              class="font-medium cursor-pointer !p-0 text-primary hover:underline"
            >
              {{ mode === "login" ? "Sign up now" : "Sign in now" }}
            </Button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
