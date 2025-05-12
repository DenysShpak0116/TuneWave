import { FC, useState } from "react"
import { AuthBtn, Container, LogoText, NavList, Wrapper } from "./header.style"
import { ROUTES } from "pages/router/consts/routes.const"
import { useNavigate } from "react-router-dom"
import { HeaderItems } from "./consts/header-item.consts"
import { NavItem } from "./nav-item"
import { useAuthStore } from "@modules/LoginForm/store/store"
import uploadIcon from "@assets/images/ic_upload.png"
import { UserBlock } from "@ui/UserHeaderBlock/user-block.component"
import search from "@assets/images/ic_search.png"
import { SearchModal } from "@modules/SearchModal"
import plusIcon from "@assets/images/ic_plus.png"

export const Header: FC = () => {
    const navigate = useNavigate()
    const { isAuth, user, logout } = useAuthStore()
    const [isSearhModalOpen, setIsSearchModalOpen] = useState<boolean>(false);

    return (
        <Wrapper>
            <Container>
                <LogoText onClick={() => navigate(ROUTES.HOME)}>TUNE WAVE</LogoText>
                <NavList>
                    <NavItem
                        title={"Пошук"}
                        icon={search}
                        onClickFn={() => setIsSearchModalOpen(true)} />

                    {HeaderItems.map((element, index) => (
                        <NavItem
                            key={index}
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

                    {isAuth() && user?.role === "admin" && (
                        <NavItem
                            path={ROUTES.ADD_CRITERION_PAGE}
                            title="Додати крітерії"
                            icon={plusIcon} />
                    )}
                </NavList>

                {isAuth() && user?.id ? (
                    <UserBlock
                        userProfileFn={() => navigate(ROUTES.USER_PROFILE.replace(":id", user?.id))}
                        profileImg={user?.profilePictureUrl}
                        username={user?.username}
                        logoutFn={logout}
                    />
                ) : (
                    <AuthBtn to={ROUTES.SIGN_IN}>Авторизуватись</AuthBtn>
                )}
            </Container>
            <SearchModal active={isSearhModalOpen} setActive={setIsSearchModalOpen} />
        </Wrapper>
    )
}
