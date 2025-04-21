import { UserInfo } from "@components/UserInfo/user-info.component";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";
import { useParams } from "react-router-dom";


export const UserProfilePage: FC = () => {
    const { id } = useParams()
    const user = useAuthStore(state => state.user)!
    const isMainUser = id === user.id
    console.log(isMainUser);


    return (
        <MainLayout>
            <UserInfo user={user} isMainUser={isMainUser} />
        </MainLayout>
    )
}