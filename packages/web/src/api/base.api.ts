import { LoginResponse } from "@modules/LoginForm/types/loginResponse";
import axios from "axios";

export const $api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    withCredentials: true
});


$api.interceptors.request.use((config) => {
    config.headers.Authorization = `Bearer ${localStorage.getItem(`token`)}`
    return config
})

$api.interceptors.response.use(
    (config) => config,
    async (error) => {
        console.log("ERROR" + error)
        const originalRequest = error.config;

        if (error.response?.status === 401 && originalRequest && !originalRequest._isRetry) {
            originalRequest._isRetry = true;
            try {
                const response = await $api.post<LoginResponse>(
                    `/auth/refresh`,
                    { withCredentials: true }
                );

                localStorage.setItem('token', response.data.accessToken);
                originalRequest.headers.Authorization = `Bearer ${response.data.accessToken}`;
                return $api.request(originalRequest);
            } catch (error: unknown) {
                console.log('Refresh failed', error);
            }
        }

        throw error;
    }
);

