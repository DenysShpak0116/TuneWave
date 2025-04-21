import { $api } from "./base.api";

export const createTrack = (formData: FormData) => {
    return $api.post("/songs", formData, {
        headers: {
            "Content-Type": "multipart/form-data",
        },
    });
};