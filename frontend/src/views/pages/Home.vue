<template>
  <div>
    <button @click="fetchUser">Get User by Firebase UID</button>
    <pre v-if="userData">{{ userData }}</pre>
    <p v-if="errorMessage" style="color: red">{{ errorMessage }}</p>
  </div>
</template>

<script lang="ts">
import UserService from "@/service/UserService"; // Assuming you have a UserService file
import AuthService from "@/service/AuthService"; // Assuming you have an AuthService file

export default {
  data() {
    return {
      userData: null as any,
      errorMessage: "",
    };
  },
  methods: {
    async fetchUser() {
      try {
        // Retrieve the Firebase access token from cookies
        const accessToken = this.getCookie("auth_token"); // Assume the cookie is called "auth_token"
        
        if (!accessToken) {
          this.errorMessage = "No access token found!";
          return;
        }

        // Authenticate user and get Firebase UID from the backend response
        const authResponse = await AuthService.authenticateUser(accessToken);

        console.log("Auth response:", authResponse);  // Log the response for debugging
      

        // Extract Firebase UID from the parsed response
        const firebaseUID = authResponse?.uid;
        
        if (!firebaseUID) {
          this.errorMessage = "No Firebase UID found in the authentication response!";
          return;
        }

        console.log("Firebase UID:", firebaseUID);

        // Fetch user data using Firebase UID
        const data = await UserService.getUserByFirebaseUID(firebaseUID, accessToken);
        console.log("Fetched data:", data);  // Log the response data
        this.userData = data;
      } catch (error: any) {
        this.errorMessage = error.message;
      }
    },
    
    // Utility method to get a cookie by name
    getCookie(name: string): string | null {
      const value = `; ${document.cookie}`;
      const parts = value.split(`; ${name}=`);
      if (parts.length === 2) return parts.pop()?.split(";").shift() || null;
      return null;
    }
  },
};
</script>
