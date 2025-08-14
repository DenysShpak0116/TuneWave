import { create } from "zustand";

export interface PlaylistItem {
    title: string;
    file: string;
    artist?: string;
    logo?: string;
    id?: string;
}

interface PlayerState {
    trackId: string;
    trackUrl: string;
    trackName: string;
    trackLogo: string;
    trackArtist: string;
    shouldAutoPlay: boolean;
    isPlaying: boolean;
    playlist: PlaylistItem[];
    setTrack: (payload: Partial<PlayerState>) => void;
    setShouldAutoPlay: (value: boolean) => void;
    setIsPlaying: (value: boolean) => void;
    setPlaylist: (playlist: PlaylistItem[] | null) => void;
    playPlayer: () => void;
    pausePlayer: () => void;
}

const defaultTrack = {
    trackId: "2af76d59-a5de-40b2-b311-bee2930b0232",
    trackUrl: "https://tunewavebucket.s3.eu-west-3.amazonaws.com/music/20b7985c-16b1-444c-a30b-3cdae7d67616/1748548464-dd.mp3",
    trackName: "Intro",
    trackLogo: "https://tunewavebucket.s3.eu-west-3.amazonaws.com/covers/20b7985c-16b1-444c-a30b-3cdae7d67616/1748548403-408c9ff55313b4d2bc89c6bea5b9224a.jpg",
    trackArtist: "TuneWave",
};

export const usePlayerStore = create<PlayerState>((set) => {
    const saved = localStorage.getItem("current-track");
    const initial = saved ? JSON.parse(saved) : defaultTrack;

    return {
        ...initial,
        shouldAutoPlay: false,
        isPlaying: false,
        playlist: [],
        setTrack: (payload) => {
            const newTrack = { ...initial, ...payload };
            localStorage.setItem("current-track", JSON.stringify(newTrack));
            set({ ...newTrack, shouldAutoPlay: true, isPlaying: true });
        },
        setShouldAutoPlay: (value) => set({ shouldAutoPlay: value }),
        setIsPlaying: (value) => set({ isPlaying: value }),
        setPlaylist: (playlist: PlaylistItem[] | null) => set({ playlist: playlist || [] }),
        pausePlayer: () => set({ isPlaying: false }),
        playPlayer: () => set({ isPlaying: true }),
    };
});
