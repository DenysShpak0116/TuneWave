import { COLORS } from "@consts/colors.consts";
import { FONTS } from "@consts/fonts.enum";
import { FC } from "react";
import { useNavigate } from "react-router-dom";
import styled from "styled-components";

interface NavItemProps {
    title: string,
    icon: string,
    path: string,
}

export const TitleContainer = styled.div`
    cursor: pointer;
    padding: 5px 20px;
    gap: 10px;
    display: flex;
    flex-direction: row;
    align-items: center;

    text-decoration: none;
    color: '${COLORS.white}';
    background-color: transparent;
    border-radius: 5px;
    transition: background-color 0.3s;

    &:hover {
        background-color: ${COLORS.dark_focusing};
    }
`

const NavIcon = styled.img`
  width: 20px;
  height: 20px;
`

const NavText = styled.p`
    font-family: ${FONTS.MONTSERRAT};
    color: ${COLORS.dark_secondary};
    font-size: 14px;
`

export const NavItem: FC<NavItemProps> = ({ title, icon, path }) => {
    const navigate = useNavigate();
    const handleClick = () => {
        navigate(path)
    }

    return (
        <>
            <TitleContainer onClick={handleClick}>
                <NavIcon src={icon} />
                <NavText>{title}</NavText>
            </TitleContainer>

        </>
    )
}