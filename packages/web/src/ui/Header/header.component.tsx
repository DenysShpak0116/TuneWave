import { FC } from "react"
import { AuthBtn, Container, LogoText, NavList, Wrapper } from "./header.style"
import { ROUTES } from "pages/router/consts/routes.const"
import { useNavigate } from "react-router-dom"
import { HeaderItems } from "./consts/header-item.consts"
import { NavItem } from "./nav-item"

export const Header: FC = () => {
    const navigate = useNavigate()

    const handleLogoClick = () => {
        navigate(ROUTES.HOME)
    }

    return (
        <>
            <Wrapper >
                <Container>
                    <LogoText onClick={handleLogoClick}>TUNE WAVE</LogoText>
                    <NavList>
                        {HeaderItems.map(element => (
                            <NavItem title={element.title} path={element.path} icon={element.icon} />
                        ))}
                    </NavList>
                    <AuthBtn to={ROUTES.SIGN_UP}>Авторизуватись</AuthBtn>
                </Container>
            </Wrapper >
        </>
    )
}