import { COLORS } from "@consts/colors.consts";
import { FONTS } from "@consts/fonts.enum";
import { Link } from "react-router-dom";
import styled from "styled-components";


export const ForgotPasswordText = styled(Link)`
    font-size: 14px;
    font-family: ${FONTS.MONTSERRAT};
    align-self: flex-start;
    color: ${COLORS.white};
`

export const RegistationText = styled.p`
    font-size: 14px;
    font-family: ${FONTS.MONTSERRAT};
    margin-top: 25px;
    cursor: pointer;

    &::before{
        content: 'Немає акаунту? ';
    }
`