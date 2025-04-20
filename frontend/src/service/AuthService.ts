import axios from 'axios';
import { API_BASE_URL } from "./apiConfig";

export default {
    async authenticateUser(accessToken: string) {
        try {
            const response = await axios.post(`${API_BASE_URL}/auth/login`, {}, {
                headers: {
                    Authorization: `Bearer ${accessToken}`, // Set the Authorization header with Bearer token
                },
                withCredentials: true, // Include cookies with the request
            });

            return response.data; // Return the response data from the server
        } catch (error) {
            console.error("Error authenticating user:", error);
            throw error; // Rethrow the error to handle it in the calling function
        }
    },
};
