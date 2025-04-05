// src/firebase/firebase.ts
import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";
import { getFirestore } from "firebase/firestore";

const firebaseConfig = {
  apiKey: "AIzaSyAzWh2jfw4ef9VeGo1Ux4ii6NJogZBgMRc",
  authDomain: "music-connect-608f6.firebaseapp.com",
  projectId: "music-connect-608f6",
  storageBucket: "music-connect-608f6.firebasestorage.app",
  messagingSenderId: "249094024240",
  appId: "1:249094024240:web:39bd9f2579018ca60e51f5",
};

const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
const db = getFirestore(app);

export { auth, db };
