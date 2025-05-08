import { ROUTES } from "pages/router/consts/routes.const"

import genre from "@assets/images/ic_genre.png"
import collections from "@assets/images/ic_collections.png"

export type NavItem = { title: string, icon: string, path: string }

export const HeaderItems: NavItem[] = [
    { title: 'Жанри', icon: genre, path: ROUTES.HOME },
    { title: "Колекції", icon: collections, path: ROUTES.HOME }
]