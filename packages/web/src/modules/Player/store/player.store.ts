import { create } from "zustand";

interface PlayerState {
    trackId: string;
    trackUrl: string;
    trackName: string;
    trackLogo: string;
    trackArtist: string;
    shouldAutoPlay: boolean;
    setTrack: (payload: Partial<PlayerState>) => void;
    setShouldAutoPlay: (value: boolean) => void;
    isPlaying: boolean;
    setIsPlaying: (value: boolean) => void;
}

const defaultTrack = {
    trackId: "5d40a144-75e1-40b6-b3d8-49c9c6733bf7",
    trackUrl: "https://tunewavebucket.s3.eu-west-3.amazonaws.com/music/20b7985c-16b1-444c-a30b-3cdae7d67616/1745244177-mySuperSong.mp3",
    trackName: "testSong",
    trackLogo: "https://tunewavebucket.s3.eu-west-3.amazonaws.com/covers/20b7985c-16b1-444c-a30b-3cdae7d67616/1745244177-ab67616d0000b273eceec97ba98bc527f6e5aec5.jpg",
    trackArtist: "sqwore, rizza",
};

export const usePlayerStore = create<PlayerState>((set) => {
    const saved = localStorage.getItem("current-track");
    const initial = saved ? JSON.parse(saved) : defaultTrack;

    return {
        ...initial,
        shouldAutoPlay: false,
        isPlaying: false,
        setTrack: (payload) => {
            const newTrack = { ...initial, ...payload };
            localStorage.setItem("current-track", JSON.stringify(newTrack));
            set({ ...newTrack, shouldAutoPlay: true, isPlaying: true });
        },
        setShouldAutoPlay: (value) => set({ shouldAutoPlay: value }),
        setIsPlaying: (value) => set({ isPlaying: value }),
    };
});
