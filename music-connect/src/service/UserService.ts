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
  async createUser(user: User) {
    const response = await axios.post<User>(USER_URL, user);
    return response.data;
  },

  async getUser(id: number) {
    const response = await axios.get<User>(`${USER_URL}/${id}`);
    return response.data;
  },

  async getAllUsers() {
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
};
