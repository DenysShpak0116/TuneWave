import { FC } from "react";
import { InputContainer, Label, StyledTextArea } from "./text-area-field.style";

interface IInputFieldProps {
    label: string;
    value: string;
    onChange: (value: string) => void;
    placeholder?: string;
}

export const TextAreaField: FC<IInputFieldProps> = ({ label, value, onChange, placeholder, }) => {
    return (
        <InputContainer>
            <Label>{label}</Label>
            <StyledTextArea
                placeholder={placeholder}
                value={value}
                onChange={(e) => onChange(e.target.value)}
            />
        </InputContainer>
    );
};