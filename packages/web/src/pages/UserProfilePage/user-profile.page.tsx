import { UserInfo } from "@components/UserInfo/user-info.component";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";
import { useParams } from "react-router-dom";
import { useGetUser } from "./hooks/useGetUserById";
import { Loader } from "@ui/Loader/loader.component";
import { SongCards } from "@components/SongCards/song-cards.component";
import { useGetUserCollections } from "./hooks/useGetUserCollections";


export const UserProfilePage: FC = () => {
    const { id } = useParams();
    const { data: user, isLoading } = useGetUser(id!);
    const mainUser = useAuthStore(state => state.user)!
    const { data: collections = [], isLoading: loadCollections } = useGetUserCollections(id!);

    const isMainUser = mainUser ? id === mainUser.id : false;

    if (isLoading || !user || loadCollections) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        );
    }

    return (
        <MainLayout>
            <UserInfo collectionsCount={collections.length} user={user} isMainUser={isMainUser} />
            <SongCards collections={collections} text="КОЛЛЕКЦІЇ КОРИСТУВАЧА" />
            <SongCards followings={user.follows.slice(0, 5)} text="ПІДПИСКИ" />
        </MainLayout>
    )
}