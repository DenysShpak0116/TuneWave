import { resetPassword } from "@api/auth.api";
import { useMutation } from "@tanstack/react-query";
import { ROUTES } from "pages/router/consts/routes.const";
import toast from "react-hot-toast";
import { useNavigate } from "react-router-dom";

export const useResetPassword = () => {
    const navigate = useNavigate();

    return useMutation({
        mutationFn: ({ newPassword, token }: { newPassword: string, token: string }) => resetPassword(newPassword, token),
        onSuccess: () => {
            toast.success("Пароль успішно змінено")
            navigate(ROUTES.SIGN_IN)
        },
        onError: (err) => {
            toast.error(`Помилка відправки коду відновлення: ${err}`)
        }
    });
};