import { createTrack } from "@api/track.api";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import toast from "react-hot-toast";
import { ErrorType } from "types/error/error.type";

export const useCreateTrack = () => {
    return useMutation({
        mutationFn: async (formData: FormData) => {
            const response = await createTrack(formData)
            return response.data;
        },
        onSuccess: () => {
            toast.success("Трек успішно створено!");
        },
        onError: (error) => {
            if (axios.isAxiosError(error) && error.response) {
                const data = error.response.data as ErrorType;
                toast.error(`Помилка створення треку: ${data.message}`);
            } else {
                toast.error("Помилка" + error.message);
            }
        }
    });
};
