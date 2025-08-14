import { COLORS } from "@consts/colors.consts";
import { FC } from "react";
import styled from "styled-components";
import googleLogo from "@assets/images/ic_google_logo.png"

const GoogleButtonContainer = styled.button`
    display: flex;
    justify-content: center;
    width: 100%;
    background-color: ${COLORS.white};
    max-width: 300px;
    border-radius: 5px;
    padding: 10px 0;
    margin-top: 28px;
`

const GoogleSpan = styled.img`
    width:20px;
    height: 20px;
`

interface GoogleButtonProps {
    onClickHandle?: VoidFunction
}

export const GoogleButton: FC<GoogleButtonProps> = ({ onClickHandle }) => {
    return (
        <GoogleButtonContainer type="button" onClick={onClickHandle}>
            <GoogleSpan src={googleLogo} />
        </GoogleButtonContainer>
    )
}