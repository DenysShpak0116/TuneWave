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

export const getTracks = async () => {
    const { data } = await $api.get("/songs");
    return data;
};

export const getTrackById = async (id: string) => {
    const { data } = await $api.get(`/songs/${id}`)
    return data
}