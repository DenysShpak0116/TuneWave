import { UserInfo } from "@components/UserInfo/user-info.component";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";
import { useParams } from "react-router-dom";
import { useGetUser } from "./hooks/useGetUserById";
import { Loader } from "@ui/Loader/loader.component";
import { useUserCollections } from "@modules/SelectCollectionModal/hooks/useUserCollections";
import { SongCards } from "@components/SongCards/song-cards.component";


export const UserProfilePage: FC = () => {
    const { id } = useParams();
    const { data: user, isLoading } = useGetUser(id!);
    const mainUser = useAuthStore(state => state.user)!
    const { data: collections = [], isLoading: loadCollections } = useUserCollections(id);

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
            <UserInfo user={user} isMainUser={isMainUser} />
            <SongCards collections={collections} text="КОЛЛЕКЦІЇ КОРИСТУВАЧА" />
        </MainLayout>
    )
}