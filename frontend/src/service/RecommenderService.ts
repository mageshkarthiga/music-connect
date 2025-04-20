import axios from "axios";
import { API_BASE_URL } from "./apiConfig";

export interface MusicRecommender {
    track_id: number;
    track_title: string;
    track_uri: string;
    track_image_url: string;
    artist_id: number;
    artist_name: string;
}

export const getRecommendedTracks = async (userId: string): Promise<MusicRecommender[]> => {
    const response = await axios.get<MusicRecommender[]>(`${API_BASE_URL}/tracks/recommendations`, {
        withCredentials: true
    });
    return response.data;
}
