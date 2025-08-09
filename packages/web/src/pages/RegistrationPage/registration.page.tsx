import { RegistrationForm } from "@modules/RegistrationForm";
import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";

export const RegistrationPage: FC = () => {

    return (
        <>
            <MainLayout>
                    <RegistrationForm />
            </MainLayout>
        </>
    )
}