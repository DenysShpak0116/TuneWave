import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const Container = styled.div`
    max-width: 400px;
    margin: 80px auto;
    padding: 32px;
    background-color: ${COLORS.dark_main};
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
    text-align: center;
`;

export const Title = styled.h2`
    font-size: 24px;
    color: ${COLORS.dark_additional};
    margin-bottom: 24px;
`;

export const Form = styled.form`
    display: flex;
    flex-direction: column;
    gap: 16px;
`;

export const Input = styled.input`
    padding: 12px;
    border-radius: 8px;
    border: 1px solid ${COLORS.dark_additional};
    font-size: 16px;
    background-color: transparent;
    color: ${COLORS.dark_additional};

    &:focus {
        outline: none;
        border: 1px solid ${COLORS.dark_focusing};
    }
`;

export const Button = styled.button`
    padding: 12px;
    border: none;
    border-radius: 8px;
    background-color: ${COLORS.dark_main};
    color: white;
    font-size: 16px;
    font-weight: bold;
    cursor: pointer;
    transition: background-color 0.3s;

    &:hover:not(:disabled) {
        background-color: ${COLORS.dark_focusing};
    }

    &:disabled {
        cursor: not-allowed;
    }
`;
