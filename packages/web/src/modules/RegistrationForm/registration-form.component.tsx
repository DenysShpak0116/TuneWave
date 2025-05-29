import { FC, useState, ChangeEvent, FormEvent } from "react";
import { LoginText } from "./registration-form.style";
import { registrationInputs } from "./consts/input.consts";
import { AuthInput } from "@ui/AuthInput/auth-input.component";
import { Button } from "@ui/Btn/btn.component";
import { GoogleButton } from "@ui/GoogleBtn/google-btn.component";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";
import { useRegister } from "./hooks/useRegistration";
import toast from "react-hot-toast";
import { AuthForm } from "@components/AuthForm/auth-form.component";

export const RegistrationForm: FC = () => {
    const navigate = useNavigate();
    const registerMutation = useRegister();

    const [formValues, setFormValues] = useState<string[]>(
        Array(registrationInputs.length).fill("")
    );

    const handleInput = (index: number) => (e: ChangeEvent<HTMLInputElement>) => {
        const newValues = [...formValues];
        newValues[index] = e.target.value;
        setFormValues(newValues);
    };

    const handleGoogleButtonClick = () => {
        try {
            const redirectUrl = import.meta.env.VITE_GOOGLE_AUTH_API
            window.location.href = redirectUrl

        } catch (e) {
            console.error(e)
        }
    };

    const handleSubmit = (e: FormEvent) => {

        e.preventDefault();

        const passwordIndex = registrationInputs.findIndex(el => el.name === "password");
        const repeatPasswordIndex = registrationInputs.findIndex(el => el.name === "repeatPassword");

        if (formValues[passwordIndex] !== formValues[repeatPasswordIndex]) {
            toast.error("Паролі не співпадають");
            return;
        }

        const requiredFields = registrationInputs.filter(el => el.name !== "repeatPassword");
        const emptyField = requiredFields.find((el, index) => {
            const value = formValues[index].trim();
            return value === "";
        });

        if (emptyField) {
            toast.error(`Поле "${emptyField.placeholder}" не може бути порожнім`);
            return;
        }

        const requestBody = requiredFields.reduce((acc, el, index) => {
            acc[el.name as "username" | "email" | "password"] = formValues[index];
            return acc;
        }, {} as { username: string; email: string; password: string });

        registerMutation.mutate(requestBody, {
            onSuccess: () => {
                toast.success("Реєстрація успішна!");
                navigate(ROUTES.SIGN_IN);
            },
            onError: () => {
                toast.error("Помилка при реєстрації. Спробуйте ще раз.");
            }
        });
    };

    return (
        <AuthForm submitFn={handleSubmit}>
            <h1>Реєстрація</h1>

            {registrationInputs.map((el, index) => (
                <AuthInput
                    key={index}
                    placeholder={el.placeholder}
                    type={el.type}
                    value={formValues[index]}
                    handleInput={handleInput(index)}
                />
            ))}

            <Button
                text={registerMutation.isPending ? "Завантаження..." : "Зареєструватися"}
                type="submit"
            />
            <GoogleButton onClickHandle={handleGoogleButtonClick} />
            <LoginText onClick={() => navigate(ROUTES.SIGN_IN)}>Авторизуйтесь</LoginText>

        </AuthForm>
    );
};
