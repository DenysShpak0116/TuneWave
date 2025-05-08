import { useMutation } from "@tanstack/react-query";
import axios from "axios";
import { useAuthStore } from "../store/store";
import toast from "react-hot-toast";
import { LoginResponse } from "../types/loginResponse";
import { ErrorType } from "types/error/error.type";
import { login } from "@api/auth.api";
import { ROUTES } from "pages/router/consts/routes.const";
import { useNavigate } from "react-router-dom";


interface LoginRequest {
    email: string;
    password: string;
}

export const useLogin = () => {
    const navigate = useNavigate()
    const setAccessToken = useAuthStore(state => state.setAccessToken);
    const setUser = useAuthStore(state => state.setUser);

    return useMutation({
        mutationFn: async (data: LoginRequest) => {
            const response = await login(data.email, data.password)
            return response.data;
        },
        onSuccess: (data: LoginResponse) => {
            setAccessToken(data.accessToken);
            setUser(data.user);
            localStorage.setItem("token", data.accessToken)
            toast.success("Вхід успішний");
            navigate(ROUTES.HOME)
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
