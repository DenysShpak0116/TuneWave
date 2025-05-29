import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const LogoContainer = styled.div`
    grid-area: "image";
    position: relative;
    width: 250px;
    height: 250px;
    border-radius: 20px;
    overflow: hidden;
`;

export const Logo = styled.img`
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 20px;
`;
export const InteractionContainer = styled.div`
    position: absolute;
    bottom: 0;
    left: 0;
    width: 100%;
    padding: 10px 0;
    display: flex;
    justify-content: space-around;
    background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent);
    border-bottom-left-radius: 20px;
    border-bottom-right-radius: 20px;
`;

export const IconButton = styled.button`
    background: none;
    border: none;
    color: white;
    cursor: pointer;
    font-size: 20px;
    display: flex;
    align-items: center;
    justify-content: center;

    &:hover {
        color: ${COLORS.dark_focusing};
    }
`;

export const InteractionIcon = styled.img`
    width: 24px;
    height: 24px;
`