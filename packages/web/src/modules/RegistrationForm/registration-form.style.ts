import { FONTS } from "@consts/fonts.enum";
import styled from "styled-components";


export const LoginText = styled.p`
    font-size: 14px;
    font-family: ${FONTS.MONTSERRAT};
    margin-top: 25px;
    cursor: pointer;

    &::before{
        content: 'Вже маєте акаунт? ';
    }
`