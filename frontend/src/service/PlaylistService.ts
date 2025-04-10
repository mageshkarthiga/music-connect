import axios from "axios";
import { API_BASE_URL } from "@/service/apiConfig";

// Define types for playlists and tracks (assuming you have the corresponding models)
export interface Playlist {
  playlist_id: number;
  playlist_name: string;
  playlist_image_url: string;
  user_id: number;
  tracks: Track[]; // Assuming a playlist has multiple tracks
}

export interface Track {
  track_id: number;
  track_title: string;
  artist_id: number;
  genre: string;
  track_uri: string;
  track_image_url: string;
  artists: Artist[];
  playlists: Playlist[];
  track_spotify_id: string;
}

export interface Artist {
  artist_id: number;
  artist_name: string;
  artist_image_url: string;
  tracks: Track[];
}

// PlaylistService class to interact with API endpoints
export default {
  // Fetch all playlists
  async getPlaylists(): Promise<Playlist[]> {
    try {
      const response = await axios.get(`${API_BASE_URL}/playlists`);
      return response.data;
    } catch (error) {
      throw new Error(`Error fetching playlists: ${error.message}`);
    }
  },

  // Fetch a single playlist by ID
  async getPlaylistById(id: number): Promise<Playlist> {
    try {
      const response = await axios.get(`${API_BASE_URL}/playlists/${id}`);
      return response.data;
    } catch (error) {
      throw new Error(
        `Error fetching playlist with ID ${id}: ${error.message}`
      );
    }
  },

  // Create a new playlist for a specific user
  async createPlaylist(
    userId: number,
    name: string,
    trackIds: number[]
  ): Promise<Playlist> {
    const playlistData = {
      user_id: userId,
      name,
      track_ids: trackIds,
    };
    try {
      const response = await axios.post(
        `${API_BASE_URL}/playlists`,
        playlistData
      );
      return response.data;
    } catch (error) {
      throw new Error(`Error creating playlist: ${error.message}`);
    }
  },

  // Update an existing playlist
  async updatePlaylist(
    id: number,
    name: string,
    trackIds: number[]
  ): Promise<Playlist> {
    const playlistData = { name, track_ids: trackIds };
    try {
      const response = await axios.put(
        `${API_BASE_URL}/playlists/${id}`,
        playlistData
      );
      return response.data;
    } catch (error) {
      throw new Error(
        `Error updating playlist with ID ${id}: ${error.message}`
      );
    }
  },

  // Delete a playlist
  async deletePlaylist(id: number): Promise<void> {
    try {
      await axios.delete(`${API_BASE_URL}/playlists/${id}`);
    } catch (error) {
      throw new Error(
        `Error deleting playlist with ID ${id}: ${error.message}`
      );
    }
  },

  // Fetch playlists for the current user
  async getPlaylistsForUser(): Promise<Playlist[]> {
    try {
      const response = await axios.get(`${API_BASE_URL}/me/playlists`, {
        withCredentials: true,
      });
      if (Array.isArray(response.data)) {
        return response.data;
      } else {
        throw new Error("Invalid playlists data format");
      }
    } catch (err: any) {
      throw new Error(`Error fetching user playlists: ${err.message}`);
    }
  },

  
  // Add a playlist for a specific user

  async addPlaylistForUser(
    name: string,    // Playlist name
    trackIds: number[],  // Array of Track IDs
    userId: number     // Single User ID (a number)
  ): Promise<Playlist> {
    const playlistData = {
      playlist_name: trackIds,  // This is the playlist name
      track_ids: userId,  // Array of track IDs
      user_id:  name,      // Single user ID
    };
  
    try {
      const response = await axios.post(
        `${API_BASE_URL}/me/playlists`,
        playlistData,
        {
          headers: {
            'Content-Type': 'application/json',
            'x-xsrf-token': getCookie('XSRF-TOKEN'),
          },
          withCredentials: true,  // Ensure cookies are sent
        }
      );
      return response.data;
    } catch (error: any) {
      if (error.response) {
        console.error("Error response from backend:", error.response);
      } else {
        console.error("Network or unknown error:", error.message);
      }
      throw new Error(
        `Error creating playlist for user: ${error.message}`
      );
    }
  },
  
  
  
  // Get playlist image URL
  async getPlaylistImage(playlistId: number): Promise<string> {
    try {
      const response = await axios.get(
        `${API_BASE_URL}/playlists/${playlistId}/image`
      );
      return response.data.image_url; // Assuming the API returns an object with image_url
    } catch (error) {
      throw new Error(
        `Error fetching playlist image for ID ${playlistId}: ${error.message}`
      );
    }
  },

  // Fetch playlists for a specific user by their userId
  async getPlaylistsByUserId(userId: number): Promise<Playlist[]> {
    try {
      const response = await axios.get(
        `${API_BASE_URL}/users/${userId}/playlists`,
        { withCredentials: true }
      );
      return response.data;
    } catch (error) {
      throw new Error(
        `Error fetching playlists for user ID ${userId}: ${error.message}`
      );
    }
  },
};

// Helper function for getting cookie (this can be moved to a separate utility file)
function getCookie(name: string): string | undefined {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop()?.split(';').shift();
  return undefined;
}
