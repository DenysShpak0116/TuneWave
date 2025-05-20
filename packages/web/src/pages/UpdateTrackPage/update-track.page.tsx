import { UpdateTrackForm } from "@modules/UpdateTrackForm";
import { MainLayout } from "@ui/layout/main-layout";
import { Loader } from "@ui/Loader/loader.component";
import { useGetTrack } from "pages/TrackPage/hooks/useGetTrack";
import { FC } from "react";
import { useParams } from "react-router-dom";

export const UpdateTrackPage: FC = () => {
    const { id } = useParams()
    const { data: track, isLoading } = useGetTrack(id!);

    if (isLoading || !track) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        );
    }

    return (
        <MainLayout>
            <UpdateTrackForm song={track} />
        </MainLayout>
    )
}