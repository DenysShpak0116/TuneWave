import { UserList } from "@modules/UserList";
import { MainLayout } from "@ui/layout/main-layout";
import { Loader } from "@ui/Loader/loader.component";
import { useGetUser } from "pages/UserProfilePage/hooks/useGetUserById";
import { FC } from "react";
import { useParams } from "react-router-dom";

export const UserListPage: FC = () => {
    const { id } = useParams()
    const { data: user, isLoading, refetch } = useGetUser(id!)

    if (isLoading) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        )
    }

    return (
        <MainLayout>
            <UserList
                refetchFn={refetch}
                followers={user!.followers}
                followings={user!.follows}
            />
        </MainLayout>
    )
}