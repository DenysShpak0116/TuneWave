import { COLORS } from "@consts/colors.consts";
import { FONTS } from "@consts/fonts.enum";
import styled from "styled-components";


export const AuthContainer = styled.div`
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
    height: calc(100vh - 100px);
`

export const FormContainer = styled.div`
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    max-width: 550px;
    border-radius: 10px;
    padding: 40px 60px;
    background-color: ${COLORS.dark_main};
`

export const LoginText = styled.p`
    font-size: 14px;
    font-family: ${FONTS.MONTSERRAT};
    margin-top: 25px;
    cursor: pointer;

    &::before{
        content: 'Вже маєте акаунт? ';
    }
`