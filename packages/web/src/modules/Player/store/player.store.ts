import { create } from "zustand";

export interface PlayerState {
    trackId: string
    trackUrl: string;
    trackName: string;
    trackLogo: string;
    trackArtist: string;
    setTrack: (payload: Omit<PlayerState, "setTrack">) => void;
}

export const usePlayerStore = create<PlayerState>((set) => ({
    trackId: "5d40a144-75e1-40b6-b3d8-49c9c6733bf7",
    trackUrl: "https://tunewavebucket.s3.eu-west-3.amazonaws.com/music/20b7985c-16b1-444c-a30b-3cdae7d67616/1745244177-mySuperSong.mp3",
    trackName: "testSong",
    trackLogo: "https://tunewavebucket.s3.eu-west-3.amazonaws.com/covers/20b7985c-16b1-444c-a30b-3cdae7d67616/1745244177-ab67616d0000b273eceec97ba98bc527f6e5aec5.jpg",
    trackArtist: "sqwore, rizza",
    setTrack: (payload) => set(() => ({ ...payload })),
}));