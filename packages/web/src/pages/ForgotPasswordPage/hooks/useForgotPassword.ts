import { sendTokenToEmail } from "@api/auth.api";
import { useMutation } from "@tanstack/react-query";
import { ROUTES } from "pages/router/consts/routes.const";
import toast from "react-hot-toast";
import { useNavigate } from "react-router-dom";

export const useSendTokenToEmail = () => {
    const navigate = useNavigate();

    return useMutation({
        mutationFn: (email: string) => sendTokenToEmail(email),
        onSuccess: () => {
            toast.success("Код відновлення відправлено на пошту")
            navigate(ROUTES.RESET_PASSWORD_PAGE)
        },
        onError: (err) => {
            toast.error(`Помилка відправки коду відновлення: ${err}`)
        }
    });
};