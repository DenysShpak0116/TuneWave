import { FC, useState } from "react"
import {
    AuthBtn,
    Burger,
    BurgerLine,
    Container,
    LogoText,
    MobileMenu,
    NavList,
    Wrapper
} from "./header.style"
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
import chatIcon from "@assets/images/ic chat.png"
import { useLogout } from "./hooks/useLogout"
import { DropdownMenu } from "@ui/DropDownMenu/drop-down-menu"
import friendsIcon from "@assets/images/ic_friends.png"

export const Header: FC = () => {
    const navigate = useNavigate()
    const { isAuth, user } = useAuthStore()
    const { mutate: logout } = useLogout()
    const [isSearhModalOpen, setIsSearchModalOpen] = useState(false)
    const [isMenuOpen, setIsMenuOpen] = useState(false)

    const handleNavigate = (fn: () => void) => {
        setIsMenuOpen(false)
        fn()
    }

    return (
        <Wrapper>
            <Container>
                <LogoText onClick={() => handleNavigate(() => navigate(ROUTES.HOME))}>
                    TUNE WAVE
                </LogoText>

                <Burger onClick={() => setIsMenuOpen(prev => !prev)}>
                    <BurgerLine />
                    <BurgerLine />
                    <BurgerLine />
                </Burger>

                <NavList>
                    <NavItem title="Пошук" icon={search} onClickFn={() => setIsSearchModalOpen(true)} />
                    {HeaderItems.map((element, index) => (
                        <NavItem key={index} title={element.title} path={element.path} icon={element.icon} />
                    ))}
                    <DropdownMenu>
                        {isAuth() && (
                            <>
                                <NavItem title="Завантажити" path={ROUTES.CREATE_TRACK} icon={uploadIcon} />
                                <NavItem title="Чати" path={ROUTES.CHAT_PAGE} icon={chatIcon} />
                                <NavItem
                                    path={ROUTES.USER_LIST.replace(":id", user!.id)}
                                    title="Підписки"
                                    icon={friendsIcon}
                                />
                            </>
                        )}
                        {isAuth() && user?.role === "admin" && (
                            <NavItem
                                path={ROUTES.ADD_CRITERION_PAGE}
                                title="Додати крітерії"
                                icon={plusIcon}
                            />
                        )}
                    </DropdownMenu>
                </NavList>

                {isAuth() && user?.id ? (
                    <UserBlock
                        userProfileFn={() => handleNavigate(() => navigate(ROUTES.USER_PROFILE.replace(":id", user.id)))}
                        profileImg={user.profilePictureUrl}
                        username={user.username}
                        logoutFn={logout}
                    />
                ) : (
                    <AuthBtn to={ROUTES.SIGN_IN}>Авторизуватись</AuthBtn>
                )}
            </Container>

            {isMenuOpen && (
                <MobileMenu>
                    <NavItem title="Пошук" icon={search} onClickFn={() => setIsSearchModalOpen(true)} />
                    {HeaderItems.map((element, index) => (
                        <NavItem
                            key={index}
                            title={element.title}
                            path={element.path}
                            icon={element.icon}
                        />
                    ))}
                    {isAuth() && (
                        <>
                            <NavItem title="Завантажити" path={ROUTES.CREATE_TRACK} icon={uploadIcon} />
                            <NavItem title="Чати" path={ROUTES.CHAT_PAGE} icon={chatIcon} />
                            <NavItem
                                path={ROUTES.USER_LIST.replace(":id", user!.id)}
                                title="Підписки"
                                icon={friendsIcon}
                            />
                        </>
                    )}
                    {isAuth() && user?.role === "admin" && (
                        <NavItem
                            path={ROUTES.ADD_CRITERION_PAGE}
                            title="Додати крітерії"
                            icon={plusIcon}
                        />
                    )}
                </MobileMenu>
            )}

            <SearchModal active={isSearhModalOpen} setActive={setIsSearchModalOpen} />
        </Wrapper>
    )
}
