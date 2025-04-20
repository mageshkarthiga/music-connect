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


export interface Friendship {
    userId: number;
    friendId: number;
    status: string; // e.g., "pending", "accepted"
    createdAt: string; // ISO date string
}


// Define the Friend interface
export interface Friend {
    count: number; 
    users: User[]; 
}

const FRIENDSHIP_URL = `${API_BASE_URL}/friends`;

// Function to send a friend request
export const sendFriendRequest = async (friendId: number) => {
    try {
        const response = await axios.post(`${FRIENDSHIP_URL}/${friendId}/request`, {
            withCredentials: true,
        });

        console.log("Friend request sent successfully:", response.data);
        return response.data; // Assuming the API returns the created friendship data

    } catch (error) {
        console.error("Error sending friend request:", error);
        throw error;
    }
};


export const getFriends = async (): Promise<Friend> => {
    try {
        const response = await axios.get(`${FRIENDSHIP_URL}`, {
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
        const response = await axios.get(`${FRIENDSHIP_URL}/pending`,
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
            `${FRIENDSHIP_URL}/${friendId}/accept`,
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
            `${FRIENDSHIP_URL}/${friendId}/reject`,
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
            `${FRIENDSHIP_URL}/${friendId}/remove`,
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



export default {
    sendFriendRequest,
    acceptFriendRequest,
    rejectFriendRequest,
    removeFriend,
    getFriends,
    getPendingFriendRequests,
};

