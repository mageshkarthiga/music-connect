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

export const getFriends = async (): Promise<Friend> => {
    try {
        const response = await axios.get(`${API_BASE_URL}/friends`, {
            withCredentials: true,
        });
        return response.data as Friend;
    } catch (error) {
        console.error("Error fetching friends:", error);
        throw error;
    }
};

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

export const acceptFriendRequest  = async (friendId: number): Promise<void> => {
    try {
        await axios.post(
            `${API_BASE_URL}/friends/${friendId}/accept`,
            {},
            {
                withCredentials: true,
            }
        );
    } catch (error) {
        console.error("Error accepting friend request:", error);
        throw error;
    }
}
export const rejectFriendRequest = async (friendId: number): Promise<void> => {
    try {
        await axios.post(
            `${API_BASE_URL}/friends/${friendId}/reject`,
            {},
            {
                withCredentials: true,
            }
        );
    } catch (error) {
        console.error("Error rejecting friend request:", error);
        throw error;
    }
}

export const removeFriend = async (friendId: number): Promise<void> => {
    try {
        await axios.post(
            `${API_BASE_URL}/friends/${friendId}/remove`,
            {},
            {
                withCredentials: true,
            }
        );
    } catch (error) {
        console.error("Error removing friend:", error);
        throw error;
    }
}