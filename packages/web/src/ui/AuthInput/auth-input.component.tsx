import { ChangeEventHandler, FC } from "react";
import { CSSProperties } from "styled-components";
import { AuthInputContainer, Input } from "./auth-input.style";


interface AuthInputProps {
    placeholder?: string;
    label?: string;
    type?: string;
    style?: CSSProperties
    value?: string | undefined;
    handleInput: ChangeEventHandler<HTMLInputElement>
}


export const AuthInput: FC<AuthInputProps> = ({ placeholder, type, style, value, handleInput }) => {
    return (
        <>
            <AuthInputContainer>
                <Input
                    placeholder={placeholder}
                    onChange={handleInput}
                    style={{ ...style }}
                    type={type || 'text'}
                    value={value}
                />
            </AuthInputContainer>
        </>
    )
}