import { CreateTrackRequest } from "../types/createTrackRequest.type";
import toast from "react-hot-toast";

export const validateCreateTrackFormData = (data: Partial<CreateTrackRequest>): data is CreateTrackRequest => {
    const requiredFields: (keyof CreateTrackRequest)[] = [
        "title", "genre", "artists", "tags", "song", "cover", "userID"
    ];

    for (const field of requiredFields) {
        const value = data[field];
        if (typeof value === "string" && value.trim() === "") {
            toast.error(`Поле "${field}" не заповнено`);
            return false;
        }
        if (Array.isArray(value) && value.length === 0) {
            toast.error(`Список "${field}" порожній`);
            return false;
        }
        if (value instanceof File && !value.name) {
            toast.error(`Файл у полі "${field}" не додано`);
            return false;
        }
        if (!value) {
            toast.error(`Поле "${field}" не заповнено`);
            return false;
        }
    }

    return true;
};
