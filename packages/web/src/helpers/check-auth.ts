import { $api } from "@api/base.api";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { LoginResponse } from "@modules/LoginForm/types/loginResponse";
export const checkAuth = async () => {
    const { setUser } = useAuthStore.getState();

    try {
        const response = await $api.get<LoginResponse>(
            "http://localhost:8081/auth/refresh",
            {
                withCredentials: true,
            }
        );
        setUser(response.data.user);
    } catch (error) {
        console.log("🔒 Auth check failed:", error);
    }
};
