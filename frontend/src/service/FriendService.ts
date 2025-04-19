import axios from "axios";
import { API_BASE_URL } from "./apiConfig"; // Make sure this points to your backend API


export interface Friendship {
    userId: number;
    friendId: number;
    status: string; // e.g., "pending", "accepted"
    createdAt: string; // ISO date string
}

export interface User {
  id?: number;
  userName: string;
  emailAddress: string;
  phoneNumber: string;
  location: string;
  profilePhotoUrl: string;
}

const FRIENDSHIP_URL = `${API_BASE_URL}/friend`; // Assuming this is the endpoint for friendships


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
// Function to accept a friend request
export const acceptFriendRequest = async (friendId: number) => {
    try {
        const response = await axios.post(`${FRIENDSHIP_URL}/${friendId}/accept`, {
            withCredentials: true,
        });
        return response.data; // Assuming the API returns the updated friendship data
    } catch (error) {
        console.error("Error accepting friend request:", error);
        throw error;
    }
};
// Function to reject a friend request
export const rejectFriendRequest = async (friendId: number) => {
    try {
        const response = await axios.post(`${FRIENDSHIP_URL}/${friendId}/reject`, {
            withCredentials: true,
        });
        return response.data; // Assuming the API returns the updated friendship data
    } catch (error) {
        console.error("Error rejecting friend request:", error);
        throw error;
    }
};
// Function to remove a friend
export const removeFriend = async (friendId: number) => {
    try {
        const response = await axios.post(`${FRIENDSHIP_URL}/${friendId}/remove`, {
            withCredentials: true,
        });
        return response.data; // Assuming the API returns the updated friendship data
    } catch (error) {
        console.error("Error removing friend:", error);
        throw error;
    }
}
// Function to get all friends
export const getFriends = async () => {
    try {
        const response = await axios.get(API_BASE_URL + "/friends", {
            withCredentials: true,
        });
        return response.data; // Assuming the API returns the list of friends
    } catch (error) {
        console.error("Error fetching friends:", error);
        throw error;
    }
}

export default {
    sendFriendRequest,
    acceptFriendRequest,
    rejectFriendRequest,
    removeFriend,
    getFriends,
};
