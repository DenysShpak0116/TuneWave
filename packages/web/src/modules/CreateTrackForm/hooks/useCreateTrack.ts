import { $api } from "@api/base.api";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import toast from "react-hot-toast";
import { ErrorType } from "types/error/error.type";

export const useCreateTrack = () => {
    return useMutation({
        mutationFn: async (formData: FormData) => {
            const response = await $api.post("/songs", formData, {
                headers: {
                    "Content-Type": "multipart/form-data"
                }
            });
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
                toast.error("Невідома помилка при створенні треку");
            }
        }
    });
};
