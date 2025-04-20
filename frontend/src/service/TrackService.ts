import axios from "axios";
import { API_BASE_URL } from "./apiConfig"; // Make sure this points to your backend API
import { supabase } from "../service/supabaseClient"; // Assuming you're using Supabase

export interface Artist {
  ArtistID: number;
  ArtistName: string;
  ArtistImageURL: string;
  Tracks: Track[]; // This holds the tracks associated with this artist
}

export interface Playlist {
  PlaylistID: number;
  PlaylistName: string;
  UserID: number;
  Tracks: Track[]; // This holds the tracks associated with this playlist
}

export interface Track {
  TrackID: number;
  TrackTitle: string;
  ArtistID: number;
  Genre: string;
  TrackURI: string;
  Artists: Artist[]; // Many-to-many relationship with Artist
  Playlists: Playlist[]; // Many-to-many relationship with Playlist
  TrackSpotifyID: string;
}

const TRACK_URL = `${API_BASE_URL}/tracks`; // Assuming this is the endpoint for tracks



// Function to fetch all tracks
export const getTracks = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/tracks`, {
      withCredentials: true,
    });
    return response.data as Track[]; // Assuming response contains a list of tracks
  } catch (error) {
    console.error("Error fetching tracks:", error);
    throw error;
  }
};

// Function to fetch a track by its ID
export const getTrackById = async (trackId: number) => {
  try {
    const response = await axios.get(`${TRACK_URL}/${trackId}`);
    return response.data as Track; // Return the track data
  } catch (error) {
    console.error("Error fetching track:", error);
    throw error;
  }
};

// Function to add a track to the user
export const addTrackForUser = async (trackData: Track) => {
  try {
    const response = await axios.post(TRACK_URL, trackData);
    return response.data; // Assuming the API returns the created track data
  } catch (error) {
    console.error("Error adding track:", error);
    throw error;
  }
};

// Function to update a track for the user
export const updateTrackForUser = async (trackId: number, trackData: Track) => {
  try {
    const response = await axios.put(`${TRACK_URL}/${trackId}`, trackData);
    return response.data; // Return updated track data
  } catch (error) {
    console.error("Error updating track:", error);
    throw error;
  }
};

// Function to delete a track for the user
export const deleteTrackForUser = async (trackId: number) => {
  try {
    const response = await axios.delete(`${TRACK_URL}/${trackId}`);
    return response.data; // Return success message or result from deletion
  } catch (error) {
    console.error("Error deleting track:", error);
    throw error;
  }
};

// Function to fetch tracks for the authenticated user based on their preferences
export const getUserTracks = async (): Promise<Track[]> => {
  try {
    const response = await axios.get(`${API_BASE_URL}/me/tracks`, {
      withCredentials: true,
    });
    return response.data as Track[];
  } catch (error) {
    console.error("Error fetching user tracks:", error);
    throw error;
  }
};

export const getFavUserTracks = async (): Promise<Track[]> => {
  try {
    const response = await axios.get(`${API_BASE_URL}/me/favtracks`, {
      withCredentials: true,
    });
    return response.data as Track[];
  } catch (error) {
    console.error("Error fetching user tracks:", error);
    throw error;
  }
};

// Function to fetch tracks for the authenticated user of another user
export const getUserTracksById = async (userId: number): Promise<Track[]> => {
  try {
    const response = await axios.get(`${API_BASE_URL}/users/${userId}/tracks`, {
      withCredentials: true,
    });
    return response.data as Track[];
  } catch (error) {
    console.error("Error fetching user tracks:", error);
    throw error;
  }
};


export const getFavUserTracksById = async (
  userId: number
): Promise<Track[]> => {
  try {
    const response = await axios.get(
      `${API_BASE_URL}/users/${userId}/favtracks`,
      {
        withCredentials: true,
      }
    );
    return response.data as Track[];
  } catch (error) {
    console.error("Error fetching user tracks:", error);
    throw error;
  }
};


// Like a track by calling the backend API
export const likeTrack = async (trackId: number) => {
  try {
    const response = await axios.post(
      `http://localhost:8080/likeTrack/${trackId}`, // Assuming this is the endpoint to like the track
      {}, // You can send any data if needed, like { is_liked: true }
      { withCredentials: true } // Include credentials (cookies) if necessary
    );

    // Handle the response if needed
    if (response.status === 201) {
      console.log("Track liked successfully:", response.data);
      return response.data;
    } else {
      console.error("Failed to like track:", response.data);
      throw new Error("Failed to like track");
    }
  } catch (error) {
    console.error("Error liking track:", error);
    throw error;
  }
}

// Unlike a track by calling the backend API
export const unlikeTrack = async (trackId: number) => {
  try {
    const response = await axios.delete(
      `http://localhost:8080/likeTrack/${trackId}`, // Assuming this is the endpoint to unlike the track
      { withCredentials: true } // Include credentials (cookies) if necessary
    );

    // Handle the response if needed
    if (response.status === 200) {
      console.log("Track unliked successfully:", response.data);
      return response.data;
    } else {
      console.error("Failed to unlike track:", response.data);
      throw new Error("Failed to unlike track");
    }
  } catch (error) {
    console.error("Error unliking track:", error);
    throw error;
  }
}

export const likedTracks = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/likedTracks`, {
      withCredentials: true,
    });
    return response.data as Track[];
  }
  catch (error) {
    console.error("Error fetching liked tracks:", error);
    throw error;
  }
};



const TrackService = {
  getTracks,
  getTrackById,
  addTrackForUser,
  updateTrackForUser,
  deleteTrackForUser,
  getUserTracks,
  getUserTracksById,
  getFavUserTracks,
  getFavUserTracksById,
  likeTrack,
  unlikeTrack,
};

export default TrackService;

  

export const incrementTrackPlayCount = async (trackId: number) => {
  try {
    const response = await axios.put(
      `${TRACK_URL}/${trackId}/increment`,
      {}, // Sending an empty body required for POST/PUT requests
      {
        withCredentials: true,
      }
    );

    return response.data.play_count;
  } catch (error) {
    console.error("Error incrementing track play count:", error);
    throw error;
  }
};
