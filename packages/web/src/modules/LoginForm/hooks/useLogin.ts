import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import { useAuthStore } from "../store/store";
import toast from "react-hot-toast";
import { LoginResponse } from "../types/loginResponse";
import { ErrorType } from "types/error/error.type";

interface LoginRequest {
    email: string;
    password: string;
}

export const useLogin = () => {
    const setAccessToken = useAuthStore(state => state.setAccessToken);
    const setUser = useAuthStore(state => state.setUser);

    return useMutation({
        mutationFn: async (data: LoginRequest) => {
            const response = await axios.post("http://localhost:8081/auth/login", data);
            return response.data;
        },
        onSuccess: (data: LoginResponse) => {
            setAccessToken(data.accessToken);
            setUser(data.user);
            localStorage.setItem("token", data.accessToken)
            toast.success("Вхід успішний");
        },
        onError: (error) => {
            if (axios.isAxiosError(error) && error.response) {
                const data = error.response.data as ErrorType;
                toast.error(`Помилка авторизації: ${data.message}`);
            } else {
                toast.error("Невідома помилка при авторизації");
            }
        }
    });
};
