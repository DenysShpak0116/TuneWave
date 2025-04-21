import { $api } from "./base.api";

export const createTrack = async (formData: FormData) => {
    return await $api.post("/songs", formData, {
        headers: {
            "Content-Type": "multipart/form-data",
        },
    });
};