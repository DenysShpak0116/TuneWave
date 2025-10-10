import { AuthForm } from "@components/AuthForm/auth-form.component";
import { FC, FormEvent, useState, ChangeEvent, useEffect } from "react";
import { loginInputs } from "./consts/input.consts";
import { AuthInput } from "@ui/AuthInput/auth-input.component";
import { Button } from "@ui/Btn/btn.component";
import { ForgotPasswordText, RegistationText } from "./login-form.style";
import { ROUTES } from "pages/router/consts/routes.const";
import { GoogleButton } from "@ui/GoogleBtn/google-btn.component";
import { useNavigate } from "react-router-dom";
import { useLogin } from "./hooks/useLogin";
import toast from "react-hot-toast";
import { encryptPassword, generateClientKeyPair } from "./utils/cryptoHelper";
import { useGetPublicKey } from "./hooks/useGetPublicKey";

export const LoginForm: FC = () => {
    const navigate = useNavigate()
    const loginMutation = useLogin();


    const [formValues, setFormValues] = useState<string[]>(
        Array(loginInputs.length).fill("")
    );


    const [clientKeys, setClientKeys] = useState<{ publicKey: string; privateKey: string } | null>(null);
    const { data: serverKey } = useGetPublicKey()

    useEffect(() => {
        (async () => {
            try {
                const keys = await generateClientKeyPair();
                setClientKeys(keys);

                localStorage.setItem("clientPrivateKey", keys.privateKey);
            } catch (err) {
                toast.error(`Can't reach security keys ${err}`);
            }
        })();
    }, []);

    const handleInput = (index: number) => (e: ChangeEvent<HTMLInputElement>) => {
        const newValues = [...formValues];
        newValues[index] = e.target.value;
        setFormValues(newValues);
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();

        const [email, password] = formValues;
        if (!email || !password) {
            toast.error("Введіть усі поля");
            return;
        }
        if (!serverKey.publicKey) {
            toast.error("Серверний ключ ще не завантажено");
            return;
        }

        try {

            const encryptedPassword = encryptPassword(serverKey.publicKey, password);

            loginMutation.mutate({ email, password: encryptedPassword, publicKey: clientKeys!.publicKey });
        } catch (err) {
            console.error("Encryption error:", err);
            toast.error("Помилка шифрування пароля");
        }
    };

    const handleGoogleButtonClick = () => {
        try {
            const redirectUrl = import.meta.env.VITE_GOOGLE_AUTH_API
            window.location.href = redirectUrl

        } catch (e) {
            console.error(e)
        }
    };

    return (
        <AuthForm submitFn={handleSubmit}>
            <h1>Авторизація</h1>

            {loginInputs.map((el, index) => (
                <AuthInput
                    key={index}
                    placeholder={el.placeholder}
                    type={el.type}
                    value={formValues[index]}
                    handleInput={handleInput(index)}
                />
            ))}

            <ForgotPasswordText to={ROUTES.FORGOT_PASSWORD_PAGE}>Забули пароль?</ForgotPasswordText>
            <Button text="Увійти" type="submit" />
            <GoogleButton onClickHandle={handleGoogleButtonClick} />
            <RegistationText onClick={() => navigate(ROUTES.SIGN_UP)}>Зареєструйтесь</RegistationText>
        </AuthForm>
    );
};
