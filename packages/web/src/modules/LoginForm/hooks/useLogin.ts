import { useMutation } from "@tanstack/react-query";
import { useAuthStore } from "../store/store";
import toast from "react-hot-toast";
import { LoginResponse } from "../types/loginResponse";
import { login } from "@api/auth.api";
import { ROUTES } from "pages/router/consts/routes.const";
import { useNavigate } from "react-router-dom";
import { LoginRequest } from "../types/loginRequest";
import { decryptWithClientPrivateKey } from "../utils/cryptoHelper";


export const useLogin = () => {
    const navigate = useNavigate();
    const setAccessToken = useAuthStore(state => state.setAccessToken);
    const setUser = useAuthStore(state => state.setUser);

    return useMutation({
        mutationFn: async (data: LoginRequest) => {
            const response = await login(data.email, data.password, data.publicKey);
            return response.data;
        },
        onSuccess: (data: LoginResponse) => {
                try {
                    const privateKey = localStorage.getItem("clientPrivateKey");
                    if (!privateKey) throw new Error("Client private key not found");

                    const decryptedToken = decryptWithClientPrivateKey(privateKey, data.accessToken);

                    setAccessToken(decryptedToken);
                    setUser(data.user);
                    localStorage.setItem("token", decryptedToken);
                    toast.success("Вхід успішний");
                    navigate(ROUTES.HOME);
                } catch (err) {
                    console.error("Decryption failed:", err);
                    toast.error("Помилка дешифрування токена");
                }
        }
    });
};
