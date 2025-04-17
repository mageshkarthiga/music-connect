import axios from "axios";
import { API_BASE_URL } from "@/service/apiConfig";

// Define types for playlists and tracks
export interface Playlist {
  playlist_id: number;
  playlist_name: string;
  playlist_image_url: string;
  user_id: number;
  tracks: Track[];
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

export interface PlaylistTrack {
  playlist_id: number;
  track_id: number;
}

// Helper function
function getCookie(name: string): string | undefined {
    const cookies = document.cookie.split('=').map(cookie => cookie.trim()); // Remove any leading/trailing spaces
    console.log("Cookies:", cookies[1]);
    return cookies[1];
  }
  
  

// PlaylistService
const PlaylistService = {
  // Fetch all playlists
  async getPlaylists(): Promise<Playlist[]> {
    try {
      const response = await axios.get(`${API_BASE_URL}/playlists`);
      return response.data;
    } catch (error: any) {
      throw new Error(`Error fetching playlists: ${error.message}`);
    }
  },

  // Fetch a single playlist by ID
  async getPlaylistById(id: number): Promise<Playlist> {
    try {
      const response = await axios.get(`${API_BASE_URL}/playlists/${id}`, {
        withCredentials: true,
      });
      return response.data;
    } catch (error: any) {
      throw new Error(`Error fetching playlist with ID ${id}: ${error.message}`);
    }
  },

  // Update an existing playlist
  async updatePlaylist(id: number, name: string, trackIds: number[]): Promise<Playlist> {
    const playlistData = { name, track_ids: trackIds };
    try {
      const response = await axios.put(
        `${API_BASE_URL}/playlists/${id}`,
        playlistData
      );
      return response.data;
    } catch (error: any) {
      throw new Error(`Error updating playlist with ID ${id}: ${error.message}`);
    }
  },

  // Delete a playlist
  async deletePlaylist(id: number): Promise<void> {
    try {
      await axios.delete(`${API_BASE_URL}/playlists/${id}`);
    } catch (error: any) {
      throw new Error(`Error deleting playlist with ID ${id}: ${error.message}`);
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

  // Create a new playlist
  async createPlaylistForUser(
    name: string,
    playlistImageUrl: string,
    userId: number
  ): Promise<Playlist> {
    const playlistData = {
      playlist_name: name,
      user_id: userId,
      playlist_image_url: playlistImageUrl,
    };

    try {
      const response = await axios.post(
        `${API_BASE_URL}/me/playlists`,
        playlistData,
        {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${getCookie('XSRF-TOKEN')}`,
          },
          withCredentials: true,
        }
      );
      return response.data;
    } catch (error: any) {
      console.error("Error creating playlist:", error.response || error.message);
      throw new Error(`Failed to create playlist: ${error.message}`);
    }
  },

  // Add tracks to an existing playlist
  async addTracksToPlaylist(
    playlistId: number,
    trackIds: number[]
  ): Promise<void> {
    const trackData = { track_ids: trackIds };

    try {
      const response = await axios.put(
        `${API_BASE_URL}/me/playlists/${playlistId}/tracks`,
        trackData,
        {
          headers: {
            'Content-Type': 'application/json',
            'x-xsrf-token': getCookie('XSRF-TOKEN'),
          },
          withCredentials: true,
        }
      );

      console.log("Tracks added response:", response.data);
    } catch (error: any) {
      console.error("Error adding tracks:", error.response || error.message);
      throw new Error(`Failed to add tracks: ${error.message}`);
    }
  },

  // Get playlist image URL
  async getPlaylistImage(playlistId: number): Promise<string> {
    try {
      const response = await axios.get(`${API_BASE_URL}/playlists/${playlistId}/image`);
      return response.data.image_url;
    } catch (error: any) {
      throw new Error(`Error fetching playlist image: ${error.message}`);
    }
  },

  // Fetch playlists for a specific user
  async getPlaylistsByUserId(userId: number): Promise<Playlist[]> {
    try {
      const response = await axios.get(
        `${API_BASE_URL}/users/${userId}/playlists`,
        { withCredentials: true }
      );
      return response.data;
    } catch (error: any) {
      throw new Error(`Error fetching playlists for user ID ${userId}: ${error.message}`);
    }
  },
};

export default PlaylistService;
