import { ROUTES } from "pages/router/consts/routes.const"
import search from "@assets/images/ic_search.png"
import genre from "@assets/images/ic_genre.png"

export type NavItem = { title: string, icon: string, path: string }

export const HeaderItems: NavItem[] = [
    { title: 'Пошук', icon: search, path: ROUTES.HOME },
    { title: 'Жанри', icon: genre, path: ROUTES.HOME }
]