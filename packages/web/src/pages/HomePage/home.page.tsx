import { SongCards } from "@components/SongCards/song-cards.component";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { MainLayout } from "@ui/layout/main-layout";
import { getUserFromCookie } from "helpers/Auth/decoder";
import { FC, useEffect, useState } from "react";
import { useTracks } from "./useTracks";
import { useGetTopCollections } from "./hooks/useGetTopCollections";
import { useUserCollections } from "@modules/SelectCollectionModal/hooks/useUserCollections";
import { UserCollections } from "@components/UserCollections/user-collections.component";

export const HomePage: FC = () => {
    const userId = useAuthStore(state => state.user?.id);
    const setAccessToken = useAuthStore((state) => state.setAccessToken);
    const setUser = useAuthStore((state) => state.setUser);
    const [limit] = useState<number>(5)
    const { data: tracks, isLoading } = useTracks({ limit: limit });
    const { data: collections, isLoading: collectionLoading } = useGetTopCollections({ limit: limit })
    const { data: userCollections = [], isLoading: loadCollections } = useUserCollections(userId!);

    useEffect(() => {
        const tryGetUser = async () => {

            const userInfo = getUserFromCookie();
            if (userInfo) {
                setAccessToken(userInfo.accessToken);
                setUser(userInfo.user);
            }
        };
        tryGetUser();
    }, []);

    return (
        <>
            <MainLayout>
                {!loadCollections && userCollections.length > 0 && (
                    <UserCollections collections={userCollections} title={"Ваші колекції"}/>
                )}
                {!isLoading && !collectionLoading && Array.isArray(tracks) && (
                    <>  
                        <SongCards songs={tracks} text="ПОПУЛЯРНІ МУЗИЧНІ ТВОРИ" />
                        <SongCards collections={collections} text="ПОПУЛЯРНІ ПЛЕЙЛИСТИ" />
                    </>

                )}
            </MainLayout>
        </>
    )
}