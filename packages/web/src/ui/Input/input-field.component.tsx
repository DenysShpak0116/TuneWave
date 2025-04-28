import { FC } from "react";
import { InputContainer, Label, StyledInput } from "./input-field.style";

interface IInputFieldProps {
    label: string;
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
    type?: string;
}

export const InputField: FC<IInputFieldProps> = ({ label, value, onChange, placeholder, type = "text" }) => {
    return (
        <InputContainer>
            <Label>{label}</Label>
            <StyledInput
                type={type}
                value={value}
                placeholder={placeholder}
                onChange={(e) => onChange(e.target.value)}
            />
        </InputContainer>
    );
};