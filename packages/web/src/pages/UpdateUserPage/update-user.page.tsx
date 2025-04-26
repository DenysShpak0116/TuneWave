import { useAuthStore } from "@modules/LoginForm/store/store"
import { UpdateUserForm } from "@modules/UpdateUserForm"
import { MainLayout } from "@ui/layout/main-layout"
import { FC } from "react"


export const UpdateUserPage: FC = () => {
    const user = useAuthStore(state => state.user!)



    return (
        <MainLayout>
            <UpdateUserForm user={user}/>
        </MainLayout>
    )
}