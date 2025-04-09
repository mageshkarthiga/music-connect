import axios from "axios";
import { API_BASE_URL } from "./apiConfig";

const EVENT_URL = `${API_BASE_URL}/events`;

export interface Event {
  event_id: number;
  event_name: string;
  event_description: string;
  event_url: string;
  event_image_url: string;
}

export default {
    async createEvent(event: Event) {
        const snakeCaseEvent = {
        event_name: event.event_name,
        event_description: event.event_description,
        event_url: event.event_url,
        event_image_url: event.event_image_url,
        };
    
        const response = await axios.post(EVENT_URL, snakeCaseEvent, {
        headers: {
            "Content-Type": "application/json",
        },
        });
        return response.data;
    },
    
    async getEvent(id: number): Promise<Event> {
        const response = await axios.get<Event>(`${EVENT_URL}/${id}`);
        return response.data;
    },
    async getAllEvents(): Promise<Event[]> {
        const response = await axios.get<Event[]>(EVENT_URL);
        return response.data;
    },
    async updateEvent(id: number, updates: Partial<Event>) {
        const response = await axios.put<Event>(`${EVENT_URL}/${id}`, updates);
        return response.data;
    },
    async deleteEvent(id: number) {
        await axios.delete(`${EVENT_URL}/${id}`);
    }
}