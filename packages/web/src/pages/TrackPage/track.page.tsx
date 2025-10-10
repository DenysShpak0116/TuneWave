import { TrackInformation } from "@modules/TrackInformation";
import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";
import { useParams } from "react-router-dom";
import { useGetTrack, useGetTrackComments } from "./hooks/useGetTrack";
import { Loader } from "@ui/Loader/loader.component";


export const TrackPage: FC = () => {
    const { id } = useParams<{ id: string }>();
    const { data: track, isLoading } = useGetTrack(id!);
    const { data: comments, isLoading: isCommentsLoading } = useGetTrackComments(id!)

    if (isLoading || isCommentsLoading || !track) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        );
    }

    return (
        <MainLayout>
            <TrackInformation song={track} songComments={comments!}/>
        </MainLayout>
    );
}