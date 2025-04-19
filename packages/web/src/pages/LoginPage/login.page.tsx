import { LoginForm } from "@modules/LoginForm";
import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";

export const LoginPage: FC = () => {

    return (
        <>
            <MainLayout>
                <LoginForm />
            </MainLayout>
        </>
    )
}