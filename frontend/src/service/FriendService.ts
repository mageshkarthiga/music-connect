import axios from "axios";
import { API_BASE_URL } from "@/service/apiConfig";

// Define the User interface
export interface User {
    userId: number;
    phoneNumber: string;
    emailAddress: string;
    location: string;
    userName: string;
    profilePhotoUrl: string;
    firebaseUid: string;
    events?: any; 
    tracks?: any; 
    friends?: any; 
}

// Define the Friend interface
export interface Friend {
    count: number; 
    users: User[]; 
}

export const getPendingFriendRequests = async (): Promise<Friend> => {
    try {
        const response = await axios.get(`${API_BASE_URL}/friends/pending`,
            {
                withCredentials: true,
            }
        );
        return response.data as Friend;
    } catch (error) {
        console.error("Error fetching pending friend requests:", error);
        throw error;
    }
};