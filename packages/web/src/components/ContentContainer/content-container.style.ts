import styled from "styled-components";
import { COLORS } from "@consts/colors.consts";

export const Container = styled.div`
    height: 100%;
    width: 100%;
    display: grid;
    grid-template-rows: 1fr auto 1fr;
    padding: 1rem;
    background-color: ${COLORS.dark_main};
    border-radius: 10px;
    color: ${COLORS.white};
    gap: 1rem;
`;

export const UploadContainer = styled.div`
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
`



export const PreviewImage = styled.img`
    max-width: 100%;
    max-height: 100%;
    border-radius: 6px;
    object-fit: cover;
`;

export const HiddenInput = styled.input`
    display: none;
`;

export const UploadBox = styled.div`
    margin-top: 0.5rem;
    height: 200px;
    border: 2px dashed ${COLORS.dark_additional};
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: border-color 0.2s ease;

    &:hover {
        border-color: ${COLORS.white};
    }
`;

export const UploadIcon = styled.div`
    color: ${COLORS.dark_additional};
    display: flex;
    align-items: center;
    justify-content: center;
    width: 200px;
    height: 200px;
    
    svg {
        width: 48px;
        height: 48px;
    }
`;


export const Divider = styled.hr`
    border: none;
    border-top: 1px solid ${COLORS.white};
`;
