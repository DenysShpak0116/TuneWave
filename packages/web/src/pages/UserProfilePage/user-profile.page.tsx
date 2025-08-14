import { FC } from "react";
import { useParams } from "react-router-dom";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { useGetUser } from "./hooks/useGetUserById";
import { useGetUserCollections } from "./hooks/useGetUserCollections";
import { MainLayout } from "@ui/layout/main-layout";
import { Loader } from "@ui/Loader/loader.component";
import { UserInfo } from "@components/UserInfo/user-info.component";
import { SongCards } from "@components/SongCards/song-cards.component";

export const UserProfilePage: FC = () => {
    const { id } = useParams<{ id: string }>();
    const mainUser = useAuthStore(state => state.user);
    const { data: user, isLoading: isUserLoading } = useGetUser(id!);
    const { data: collections = [], isLoading: isCollectionsLoading } = useGetUserCollections(id!);
    const isMainUser = mainUser?.id === id;
    const isLoading = isUserLoading || isCollectionsLoading || !user;

    if (isLoading) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        );
    }

    const userCollections = collections.filter(c => c.user?.id === id);
    const savedCollections = collections.filter(c => c.user?.id !== id);
    console.log(savedCollections)

    return (
        <MainLayout>
            <UserInfo
                user={user}
                isMainUser={isMainUser}
                collectionsCount={userCollections.length}
            />
            {userCollections.length > 0 && (
                <SongCards collections={userCollections} text="КОЛЛЕКЦІЇ КОРИСТУВАЧА" />
            )}
            {savedCollections.length > 0 && (
                <SongCards collections={savedCollections} text="ЗБЕРЕЖЕНІ КОЛЕКЦІЇ" />
            )}
            {user.follows.length > 0 && (
                <SongCards followings={user.follows.slice(0, 5)} text="ПІДПИСКИ" />
            )}
        </MainLayout>
    );
};
