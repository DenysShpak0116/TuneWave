import { COLORS } from "@consts/colors.consts";
import { FONTS } from "@consts/fonts.enum";
import styled from "styled-components";


export const AuthInputContainer = styled.div`
    display: flex;
    flex-direction: column;
    width: 100%;
    padding: 20px 113px;
`

export const Input = styled.input`
    width: 100%;
    border-radius: 5px;
    font-family: ${FONTS.MONTSERRAT};
    font-size: 16px;
    font-weight: 400;
    background-color: transparent;
    border: 2px solid ${COLORS.white};
    height: 40px;
    padding-left: 10px;
    color: ${COLORS.white}
`