import axios from "axios";
import { API_BASE_URL } from "./apiConfig";
import { supabase } from "../service/supabaseClient";
import { auth } from "@/firebase/firebase";
import { getUserTracksById } from "./TrackService";

export interface User {
  id?: number;
  userName: string;
  emailAddress: string;
  phoneNumber: string;
  location: string;
  profilePhotoUrl: string;
}

export default {
  async createUser(user: User & { firebaseUID: string }) {
    const { data, error } = await supabase
      .from("users")
      .insert([
        {
          user_name: user.userName,
          phone_number: user.phoneNumber,
          email_address: user.emailAddress,
          location: user.location,
          profile_photo_url: user.profilePhotoUrl,
          firebase_uid: user.firebaseUID,
        },
      ])
      .single(); // Use `.single()` to get a single response

    if (error) {
      console.error("Error creating user:", error.message);
      throw error; // or return null/error
    }

    return data; // Return the inserted user data
  },

  async getUser() {
    const response = await axios.get<User>(`${API_BASE_URL}/me`, {
      withCredentials: true,
    });
    return response.data;
  },

  async getAllUsers(): Promise<User[]> {
    const response = await axios.get<User[]>(`${API_BASE_URL}/users`, {
      withCredentials: true,
    });
    return response.data;
  },

  async updateUser(id: number, updates: Partial<User>) {
    const response = await axios.put<User>(`${API_BASE_URL}/${id}`, updates);
    return response.data;
  },

  async deleteUser(id: number) {
    await axios.delete(`${API_BASE_URL}/${id}`);
  },

  async getUserByFirebaseUID(firebaseUID: string): Promise<User> {
    const response = await axios.get(`${API_BASE_URL}/users/firebase/${firebaseUID} `, {
      withCredentials: true,
    });
    const data = response.data;

    return data;
  },

  async fetchSecureData(idToken: string) {
    try {
      const response = await axios.get(`${API_BASE_URL}/secure`, {
        headers: {
          Authorization: `Bearer ${idToken}`,
        },
      });
      return response.data;
    } catch (error) {
      console.error("Error fetching secure data", error);
      throw error;
    }
  },

  async getUserByUserId(userId: number): Promise<User> {
    const response = await axios.get(`${API_BASE_URL}/users/${userId}`, {
      withCredentials: true,
    });

    return response.data;
  },

  
};
