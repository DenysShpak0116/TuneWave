import { createCollection } from "@api/collection.api";
import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import toast from "react-hot-toast";
import { ErrorType } from "types/error/error.type";

export const useCreateCollection = () => {
    return useMutation({
        mutationFn: async (formData: FormData) => {
            const data = await createCollection(formData);
            return data;
        },
        onError: (error) => {
            if (axios.isAxiosError(error) && error.response) {
                const data = error.response.data as ErrorType;
                toast.error(`Помилка створення колекції: ${data.message}`);
            } else {
                console.error("Error: ", error);
                toast.error("Сталася помилка під час створення колекції.");
            }
        }
    });
};