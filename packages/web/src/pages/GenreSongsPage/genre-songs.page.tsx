import { SongCards } from "@components/SongCards/song-cards.component";
import { MainLayout } from "@ui/layout/main-layout";
import { Loader } from "@ui/Loader/loader.component";
import { useTracks } from "pages/HomePage/useTracks";
import { FC } from "react";
import { useParams } from "react-router-dom";


export const GenreSongsPage: FC = () => {
    const { genre } = useParams()
    const { data: tracks, isLoading } = useTracks({ search: genre });

    if (isLoading) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        )
    }
    return (
        <MainLayout>
            <SongCards songs={tracks} text={`Пісні за жанром ${genre}`}/>
        </MainLayout>
    )
}