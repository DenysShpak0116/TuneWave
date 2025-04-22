import { SongCards } from "@components/SongCards/song-cards.component";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { MainLayout } from "@ui/layout/main-layout";
import { getUserFromCookie } from "helpers/Auth/decoder";
import { FC, useEffect } from "react";
import { useTracks } from "./useTracks";

export const HomePage: FC = () => {
    const setAccessToken = useAuthStore((state) => state.setAccessToken);
    const setUser = useAuthStore((state) => state.setUser);
    const { data: tracks, isLoading } = useTracks();

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
                {!isLoading && Array.isArray(tracks) && (
                    <SongCards songs={tracks} />
                )}
            </MainLayout>
        </>
    )
}