import { COLORS } from "@consts/colors.consts";
import { FONTS } from "@consts/fonts.enum";
import { Link } from "react-router-dom";
import styled from "styled-components";

export const Wrapper = styled.div`
    background-color: ${COLORS.dark_main};
    margin-bottom: 20px;
`

export const Container = styled.div`
    max-width: 1280px;
    width: 100%;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    padding: 12px 24px;
    margin: 0 auto;
`

export const LogoText = styled.div`
    font-family: ${FONTS.JERSEY25};
    font-size: 36px;
    justify-content: flex-start;
    color: ${COLORS.white};
    cursor: pointer;
`

export const NavList = styled.div`
    display: flex;
    flex-direction: row;
`

export const AuthBtn = styled(Link)`
    justify-self: end;
    padding: 10px 20px;
    text-decoration: none;
    color: ${COLORS.dark_secondary};
    background-color: transparent;
    border-radius: 5px;
    transition: background-color 0.3s;

    &:hover {
        background-color: ${COLORS.dark_focusing};
    }
`