import { GetTracksParams } from "types/song/getTrackParams";
import { $api } from "./base.api";

export const createTrack = async (formData: FormData) => {
    try {
        const { data } = await $api.post("/songs", formData, {
            headers: {
                "Content-Type": "multipart/form-data",
            },
        })
        return data
    }
    catch (err) {
        console.log(err)
    }
};

export const updateTrack = async (songId: string, formData: FormData) => {
    try {
        const { data } = await $api.put(`/songs/${songId}`, formData, {
            headers: {
                "Content-Type": "multipart/form-data",
            },
        })
        return data
    }
    catch (err) {
        console.log(err)
    }
};

export const getTracks = async (params: GetTracksParams = {}) => {
    const { data } = await $api.get("/songs", { params });
    return data;
};

export const deleteTrack = async (id: string) => {
    const { data } = await $api.delete(`/songs/${id}`)
    return data
}

export const getTrackById = async (id: string) => {
    const { data } = await $api.get(`/songs/${id}`)
    return data
}

export const reactSong = async (songId: string, reactionType: "like" | "dislike", userId: string) => {
    const { data } = await $api.post(`/songs/${songId}/reaction`, { reactionType, userId })
    return data;
}

export const getUserReaction = async (songId: string, userId?: string) => {
    const { data } = await $api.get(`/songs/${songId}/is-reacted/${userId}`)
    return data;
}

export const addSongToCollection = async (trackId: string, collectionId: string) => {
    const { data } = await $api.post(`/songs/${trackId}/add-to-collection`, { collectionId: collectionId })
    return data
}