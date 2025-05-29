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
  gap: 10px;

  @media (max-width: 768px) {
    display: none;
  }
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

export const Burger = styled.div`
  display: none;
  flex-direction: column;
  justify-content: center;
  gap: 4px;
  cursor: pointer;
  margin-left: 16px;

  @media (max-width: 768px) {
    display: flex;
  }
`

export const BurgerLine = styled.span`
  width: 25px;
  height: 3px;
  background-color: ${COLORS.white};
  border-radius: 2px;
`

export const MobileMenu = styled.div`
  display: flex;
  flex-direction: column;
  background-color: ${COLORS.dark_backdrop};
  padding: 16px 24px;
  gap: 12px;

  @media (min-width: 769px) {
    display: none;
  }

  a {
    text-decoration: none;
    color: ${COLORS.white};
    font-family: ${FONTS.MONTSERRAT};
    font-size: 16px;

    &:hover {
      color: ${COLORS.dark_focusing};
    }
  }
`

