import styled from "styled-components";
import { COLORS } from "@consts/colors.consts";


export const SongDetailsContainer = styled.div`
    height: 100%;
    width: 100%;
    background-color: ${COLORS.dark_main};
    border-radius: 10px;
    color: ${COLORS.white};
    display: flex;
    align-items: center;
    flex-direction: column;
    gap: 1rem;
    padding: 1rem;
`

export const InputsWrapper = styled.div`
    display: flex;
    flex-direction: column;
    gap: 16px;
    width: 100%;
`;

export const InputGroup = styled.div`
    display: flex;
    flex-direction: column;
    gap: 10px;

    label {
        font-size: 14px;
        color: ${COLORS.white};
    }
`;

export const StyledInput = styled.input`
    padding: 10px;
    border: 1px solid ${COLORS.dark_additional};
    border-radius: 6px;
    background-color: transparent;
    color: ${COLORS.white};

    &:focus {
        outline: none;
        border-color: ${COLORS.white};
    }
`;

export const StyledTextarea = styled.textarea`
    padding: 10px;
    border: 1px solid ${COLORS.dark_additional};
    border-radius: 6px;
    background-color: transparent;
    color: ${COLORS.white};
    resize: vertical;

    &:focus {
        outline: none;
        border-color: ${COLORS.white};
    }
`;
