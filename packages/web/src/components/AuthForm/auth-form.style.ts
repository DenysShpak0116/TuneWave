import { COLORS } from "@consts/colors.consts"
import styled from "styled-components"

export const AuthContainer = styled.div`
    flex: 1;
    display: flex;
    justify-content: center;
    align-items: center;
    height: calc(100vh - 100px);
`

export const FormContainer = styled.form`
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    max-width: 550px;
    border-radius: 10px;
    padding: 40px 60px;
    background-color: ${COLORS.dark_main};
`