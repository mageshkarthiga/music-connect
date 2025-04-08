// EventService.ts

import axios from "axios";
import { API_BASE_URL } from "./apiConfig";
import { supabase } from "../service/supabaseClient";
import Cookies from "js-cookie"; // Import js-cookie to access cookies

export interface Event {
  id?: number;
  eventName: string;
  eventDescription: string;
  eventURL: string;
  eventImageURL: string;
}

const EVENT_URL = `${API_BASE_URL}/events`;

export default {
  // Supabase creation only
  async createEvent(event: Event) {
    const { data, error } = await supabase
      .from("events")
      .insert([
        {
          event_name: event.eventName,
          event_description: event.eventDescription,
          event_url: event.eventURL,
          event_image_url: event.eventImageURL,
        },
      ])
      .select()
      .single();

    if (error) {
      console.error("Error creating event in Supabase:", error.message);
      throw error;
    }

    return data;
  },

  async getEvent(id: number) {
    const response = await axios.get<Event>(`${EVENT_URL}/${id}`);
    return response.data;
  },

  async getAllEvents() {
    const response = await axios.get<Event[]>(EVENT_URL);
    return response.data;
  },

  async updateEvent(id: number, updates: Partial<Event>) {
    const response = await axios.put<Event>(`${EVENT_URL}/${id}`, updates);
    return response.data;
  },

  async deleteEvent(id: number) {
    await axios.delete(`${EVENT_URL}/${id}`);
  },

  async getEventsByUserId(userId: number) {
    // Extract auth token from cookies
    const authToken = Cookies.get("auth_token"); // Access the token from the cookies
    
    if (!authToken) {
      throw new Error("Auth token not found.");
    }

    const response = await axios.get<Event[]>(`${API_BASE_URL}/users/${userId}/events`, {
      headers: {
        Authorization: `Bearer ${authToken}`, // Use the extracted token here
      },
      withCredentials: true, // Ensure cookies are sent automatically with the request
    });
  
    return response.data;
  
}

  
  

//   async getUserByEventId(eventId: number, token: string) {
//     const response = await axios.get<Event>(`${EVENT_URL}/users/${eventId}`, {
//       headers: {
//         Authorization: `Bearer ${token}`,
//       },
//     });
//     return response.data;
  
  
};
