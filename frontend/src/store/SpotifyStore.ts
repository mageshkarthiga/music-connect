import { defineStore } from "pinia";

export const useSpotifyStore = defineStore("spotify", {
    state: () => ({
        spotifyUri: "spotify:track:3lzUeaCbcCDB5IXYfqWRlF" as string, // The current Spotify track URI
    }),
    actions: {
        setSpotifyUri(uri) {
            this.spotifyUri = uri;
        },
    },
});