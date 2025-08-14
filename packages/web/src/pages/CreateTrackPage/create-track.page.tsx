import { CreateTrackForm } from "@modules/CreateTrackForm";
import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";



export const CreateTrackPage: FC = () => {

    return (
        <MainLayout>
            <CreateTrackForm />
        </MainLayout>
    )
}