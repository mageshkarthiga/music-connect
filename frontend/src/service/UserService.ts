import axios from "axios";
import { API_BASE_URL } from "./apiConfig";

const USER_URL = `${API_BASE_URL}/users`;

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
    const snakeCaseUser = {
      user_name: user.userName,
      phone_number: user.phoneNumber,
      email_address: user.emailAddress,
      location: user.location,
      profile_photo_url: user.profilePhotoUrl,
      firebase_uid: user.firebaseUID,
    };

    const response = await axios.post(USER_URL, snakeCaseUser, {
      headers: {
        "Content-Type": "application/json",
      },
    });
    return response.data;
  },

  async getUser(id: number): Promise<User> {
    const response = await axios.get<User>(`${USER_URL}/${id}`);
    return response.data;
  },

  async getAllUsers(): Promise<User[]> {
    const response = await axios.get<User[]>(USER_URL);
    return response.data;
  },

  async updateUser(id: number, updates: Partial<User>) {
    const response = await axios.put<User>(`${USER_URL}/${id}`, updates);
    return response.data;
  },

  async deleteUser(id: number) {
    await axios.delete(`${USER_URL}/${id}`);
  },

  async getUserByFirebaseUID(firebaseUID: string): Promise<User> {
    const response = await axios.get(`${USER_URL}/firebase/${firebaseUID}`);
    const data = response.data;

    return {
      id: data.user_id,
      userName: data.user_name,
      emailAddress: data.email_address,
      phoneNumber: data.phone_number,
      location: data.location,
      profilePhotoUrl: data.profile_photo_url,
    };
  },
};
