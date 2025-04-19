
import { useAuthStore } from "@modules/LoginForm/store/store";
import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";

export const HomePage: FC = () => {

    const user = useAuthStore((state) => state.user);

    return (
        <>
            <MainLayout>
                <h1>Hello {user?.username}!</h1>
            </MainLayout>
        </>
    )
}