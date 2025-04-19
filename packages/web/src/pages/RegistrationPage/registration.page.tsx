
import { RegistrationForm } from "@modules/RegistrationForm/registration-form.component";
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