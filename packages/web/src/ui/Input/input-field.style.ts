import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const InputContainer = styled.div`
    display: flex;
    flex-direction: column;
    gap: 8px;
`;

export const Label = styled.label`
    font-size: 14px;
    color: ${COLORS.dark_additional};
    margin-top: 10px;
`;

export const StyledInput = styled.input`
    width: 100%;
    padding: 10px 0;
    background-color: transparent;
    border: 2px solid ${COLORS.dark_backdrop};
    border-radius: 6px;
    color: white;
    font-size: 16px;
    text-indent: 5px;

    &:focus {
        outline: none;
        border-color: ${COLORS.dark_focusing};
    }
`;