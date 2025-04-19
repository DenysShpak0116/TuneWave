import { FC } from "react";
import { StyledButton } from "./btn.style";
import { CSSProperties } from "styled-components";

interface ButtonProps {
    style?: CSSProperties;
    text: string;
    onClick?: VoidFunction
    type?: "button" | "submit" | "reset";
}

export const Button: FC<ButtonProps> = ({ style, text, onClick, type = "button" }) => {
    return (
        <StyledButton style={style} onClick={onClick} type={type}>
            {text}
        </StyledButton>
    );
};