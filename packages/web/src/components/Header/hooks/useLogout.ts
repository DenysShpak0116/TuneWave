import { logout } from "@api/auth.api";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { useMutation } from "@tanstack/react-query";
import { ROUTES } from "pages/router/consts/routes.const";
import { useNavigate } from "react-router-dom";

export const useLogout = () => {
    const navigate = useNavigate()
    const {logout : authLogout} = useAuthStore()

    return useMutation({
        mutationFn: async () => {
            await logout()
            authLogout()
        },
        onSuccess: () => {
            navigate(ROUTES.HOME)
        },
    });
};