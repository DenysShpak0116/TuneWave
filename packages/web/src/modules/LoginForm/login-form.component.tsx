import { AuthForm } from "@components/AuthForm/auth-form.component";
import { FC, FormEvent, useState, ChangeEvent } from "react";
import { loginInputs } from "./consts/input.consts";
import { AuthInput } from "@ui/AuthInput/auth-input.component";
import { Button } from "@ui/Btn/btn.component";
import { ForgotPasswordText, RegistationText } from "./login-form.style";
import { ROUTES } from "pages/router/consts/routes.const";
import { GoogleButton } from "@ui/GoogleBtn/google-btn.component";
import { useNavigate } from "react-router-dom";
import { useLogin } from "./hooks/useLogin";
import toast from "react-hot-toast";

export const LoginForm: FC = () => {
    const navigate = useNavigate()
    const loginMutation = useLogin();

    const [formValues, setFormValues] = useState<string[]>(
        Array(loginInputs.length).fill("")
    );

    const handleInput = (index: number) => (e: ChangeEvent<HTMLInputElement>) => {
        const newValues = [...formValues];
        newValues[index] = e.target.value;
        setFormValues(newValues);
    };

    const handleSubmit = (e: FormEvent) => {
        e.preventDefault();

        const [email, password] = formValues

        if (email == "" || password == "") {
            toast.error('Введіть усі поля')
            return
        }

        loginMutation.mutate({ email, password });
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

            <ForgotPasswordText to={ROUTES.HOME}>Забули пароль?</ForgotPasswordText>
            <Button text="Увійти" type="submit" />
            <GoogleButton onClickHandle={handleGoogleButtonClick} />
            <RegistationText onClick={() => navigate(ROUTES.SIGN_UP)}>Зареєструйтесь</RegistationText>
        </AuthForm>
    );
};
