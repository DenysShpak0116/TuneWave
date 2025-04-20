import { useAuthStore } from "@modules/LoginForm/store/store";
import { MainLayout } from "@ui/layout/main-layout";
import { getUserFromCookie } from "helpers/Auth/decoder";
import { FC, useEffect } from "react";

export const HomePage: FC = () => {
    const setAccessToken = useAuthStore((state) => state.setAccessToken);
    const setUser = useAuthStore((state) => state.setUser);
    
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

    const user = useAuthStore((state) => state.user);

    return (
        <>
            <MainLayout>
                <h1>Hello {user?.username}!</h1>
            </MainLayout>
        </>
    )
}