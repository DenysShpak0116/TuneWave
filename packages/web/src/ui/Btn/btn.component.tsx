import { FC } from "react";
import { StyledButton } from "./btn.style";
import { CSSProperties } from "styled-components";

interface ButtonProps {
    style?: CSSProperties;
    text: string;
    onClick?: VoidFunction
    type?: "button" | "submit" | "reset";
    isDisabled?: boolean
}

export const Button: FC<ButtonProps> = ({ style, text, onClick, type = "button", isDisabled }) => {
    return (
        <StyledButton disabled={isDisabled} style={style} onClick={onClick} type={type}>
            {text}
        </StyledButton>
    );
};