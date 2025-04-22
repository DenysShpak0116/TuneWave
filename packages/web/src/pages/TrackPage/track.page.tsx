import { TrackInformation } from "@modules/TrackInformation";
import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";
import { useParams } from "react-router-dom";
import { useGetTrack } from "./hooks/useGetTrack";
import { Loader } from "@ui/Loader/loader.component";


export const TrackPage: FC = () => {
    const { id } = useParams<{ id: string }>();
    const { data: track, isLoading } = useGetTrack(id!);

    if (isLoading) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        );
    }

    return (
        <MainLayout>
            <TrackInformation song={track} />
        </MainLayout>
    );
}