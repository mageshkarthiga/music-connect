// EventService.ts

import axios from "axios";
import { API_BASE_URL } from "./apiConfig";
import { supabase } from "../service/supabaseClient";

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
    const response = await axios.get<Event[]>(EVENT_URL, {
      withCredentials: true,
    });
    return response.data;
  },

  async getAllEventVenues() {
    const response = await axios.get<Event[]>(`${EVENT_URL}/venues`, {
      withCredentials: true,
    });
    return response.data;
  },

  async updateEvent(id: number, updates: Partial<Event>) {
    const response = await axios.put<Event>(`${EVENT_URL}/${id}`, updates);
    return response.data;
  },

  async deleteEvent(id: number) {
    await axios.delete(`${EVENT_URL}/${id}`);
  },
  async getEventsForCurrentUser(): Promise<Event[]> {
    const response = await axios.get<Event[]>(`${API_BASE_URL}/me/events`, {
      withCredentials: true, // Ensures cookies (like Firebase auth token) are sent
    });

    return response.data;
  },
  async getFavEventsForCurrentUser(): Promise<Event[]> {
    const response = await axios.get<Event[]>(`${API_BASE_URL}/me/favevents`, {
      withCredentials: true, // Ensures cookies (like Firebase auth token) are sent
    });

    return response.data;
  },
  async getEventsByUserId(userId: number): Promise<Event> {
    const response = await axios.get<Event>(
      `${API_BASE_URL}/users/${userId}/events`,
      {
        withCredentials: true,
      }
    );
    return response.data;
  },

  async getFavEventsByUserId(userId: number): Promise<Event> {
    const response = await axios.get<Event>(
      `${API_BASE_URL}/users/${userId}/favevents`,
      {
        withCredentials: true,
      }
    );
    return response.data;
  },

  async likeEvent(eventId) {
   
    return await axios.post(
      `${API_BASE_URL}/likeEvent`,
      { event_id: eventId }, // ✅ Must match expected format
      {
        withCredentials: true,
      }
    );
  },
  

  async getLikedEvents(): Promise<Event[]> {
 
    const response = await axios.get(`${API_BASE_URL}/me/likedEvents`, {
      withCredentials: true,
    });
    return response.data;
  },

// EventService.js or inside your component
async unlikeEvent(eventId) {
  try {
    const response = await axios.delete(`${API_BASE_URL}/likeEvent`, {
      withCredentials: true,
      data: { event_id: eventId },  // Send the event_id in the request body
    });
    return response.data;
  } catch (error) {
    console.error("Failed to unlike event:", error);
    throw error;
  }
}

};
