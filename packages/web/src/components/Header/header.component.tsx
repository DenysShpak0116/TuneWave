import { FC } from "react"
import { AuthBtn, Container, LogoText, NavList, Wrapper } from "./header.style"
import { ROUTES } from "pages/router/consts/routes.const"
import { useNavigate } from "react-router-dom"
import { HeaderItems } from "./consts/header-item.consts"
import { NavItem } from "./nav-item"
import { useAuthStore } from "@modules/LoginForm/store/store"
import uploadIcon from "@assets/images/ic_upload.png"
import { UserBlock } from "@ui/UserHeaderBlock/user-block.component"

export const Header: FC = () => {
    const navigate = useNavigate()
    const { isAuth, user, logout } = useAuthStore()

    const handleLogoClick = () => {
        navigate(ROUTES.HOME)
    }

    return (
        <Wrapper>
            <Container>
                <LogoText onClick={handleLogoClick}>TUNE WAVE</LogoText>
                <NavList>
                    {HeaderItems.map((element) => (
                        <NavItem
                            key={element.path}
                            title={element.title}
                            path={element.path}
                            icon={element.icon}
                        />
                    ))}

                    {isAuth() && (
                        <NavItem
                            title="Завантажити"
                            path={ROUTES.CREATE_TRACK}
                            icon={uploadIcon}
                        />
                    )}
                </NavList>

                {isAuth() ? (
                    <UserBlock
                        profileImg={user?.profilePictureUrl}
                        username={user?.username}
                        logoutFn={logout}
                    />
                ) : (
                    <AuthBtn to={ROUTES.SIGN_IN}>Авторизуватись</AuthBtn>
                )}
            </Container>
        </Wrapper>
    )
}
