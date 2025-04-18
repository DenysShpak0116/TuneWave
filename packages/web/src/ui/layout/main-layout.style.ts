import { COLORS } from "@consts/colors.consts";
import styled from "styled-components";


export const Wrapper = styled.div`
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    background-color: ${COLORS.dark_backdrop};
    color: ${COLORS.dark_additional};
`

export const Container = styled.div`
    max-width: 1280px;
    width: 100%;
    margin: 0 auto;
    gap: 24px
`

