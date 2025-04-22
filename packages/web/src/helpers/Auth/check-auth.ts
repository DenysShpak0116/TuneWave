import { $api } from "@api/base.api";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { LoginResponse } from "@modules/LoginForm/types/loginResponse";

export const checkAuth = async () => {
    const { setUser, setAccessToken } = useAuthStore.getState();

    const cookie = document.cookie
        .split("; ")
        .find((row) => row.startsWith("refreshToken="));
    if (!cookie) {
        console.warn("refreshToken cookie not found");
        return;
    }
    const refreshToken = cookie.split("=")[1];

    try {
        const response = await $api.post<LoginResponse>(
            "/auth/refresh",
            { refreshToken },
            { withCredentials: true }
        );

        setAccessToken(response.data.accessToken);
        setUser(response.data.user);
    } catch (error) {
        console.error("Failed to refresh token:", error);
    }
};