import { MainLayout } from "@ui/layout/main-layout";
import { FC } from "react";
import { useGetGenres } from "./hooks/useGetGenres";
import { Loader } from "@ui/Loader/loader.component";
import { GenreList } from "@components/GenreList/genre-list.component";

export const GenrePage: FC = () => {
    const {data: genres, isLoading} = useGetGenres()

    if(isLoading){
        return(
            <MainLayout>
                <Loader/>
            </MainLayout>
        )
    }


    return(
        <MainLayout>
            <GenreList genres={genres!}/>
        </MainLayout>
    )
}