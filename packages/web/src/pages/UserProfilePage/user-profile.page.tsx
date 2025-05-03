import { UserInfo } from "@components/UserInfo/user-info.component";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";
import { useParams } from "react-router-dom";
import { useGetUser } from "./hooks/useGetUserById";
import { Loader } from "@ui/Loader/loader.component";


export const UserProfilePage: FC = () => {
    const { id } = useParams();
    const { data: user, isLoading } = useGetUser(id!);
    const mainUser = useAuthStore(state => state.user)!

    const isMainUser = mainUser ? id === mainUser.id : false;

    if (isLoading || !user) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        );
    }

    return (
        <MainLayout>
            <UserInfo user={user} isMainUser={isMainUser} />
        </MainLayout>
    )
}