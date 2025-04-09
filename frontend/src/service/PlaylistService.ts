import axios from "axios";
import { API_BASE_URL } from "@/service/apiConfig";

// Define types for playlists and tracks (assuming you have the corresponding models)
export interface Track {
  id: number;
  name: string;
}

export interface Playlist {
  id: number;
  name: string;
  user_id: number;
  playlist_image_url: string;
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

  // Create a new playlist
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

  // Fetch playlists for a specific user
  async getPlaylistsForUser(): Promise<Playlist[]> {
    try {
      const response = await axios.get(`${API_BASE_URL}/me/playlists`, {
        withCredentials: true,
      });
      console.log("Fetched playlists:", response.data);
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
    userId: number,
    name: string,
    trackIds: number[]
  ): Promise<Playlist> {
    const playlistData = {
      name,
      track_ids: trackIds,
    };
    try {
      const response = await axios.post(
        `${API_BASE_URL}/playlists/user/${userId}`,
        playlistData
      );
      return response.data;
    } catch (error) {
      throw new Error(
        `Error creating playlist for user ${userId}: ${error.message}`
      );
    }
  },
  //get playlist image 

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
};
