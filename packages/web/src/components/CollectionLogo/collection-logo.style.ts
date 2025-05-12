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