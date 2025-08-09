import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";

export const StyledPlayerContainer = styled.div`
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 60px;
    background-color: ${COLORS.dark_main};
    z-index: 999;
    display: flex;
    align-items: center;
    justify-content: space-between;
`;

export const PlayerInsideContainer = styled.div`
    max-width: 1280px;
    margin: 0 auto;
    width: 100%;
    display: flex;
    justify-content: space-around;

`