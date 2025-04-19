import { FC, FormEvent, ReactNode } from "react";
import { AuthContainer, FormContainer } from "./auth-form.style";

interface AuthFormProps {
    submitFn: (e: FormEvent) => void,
    children: ReactNode

}

export const AuthForm: FC<AuthFormProps> = ({ submitFn, children }) => {

    return (
        <AuthContainer>
            <FormContainer onSubmit={submitFn}>
                {children}
            </FormContainer>
        </AuthContainer>
    )
}