import { FC } from "react";
import { GenreResponseType } from "types/song/genreResponseType";
import { Card, CoverImage, CoverWrapper, GenreTitle, Grid } from "./genre-list.component.style";

interface IGenreListProps {
    genres: GenreResponseType[];
}

export const GenreList: FC<IGenreListProps> = ({ genres }) => {
    return (
        <Grid>
            {genres.map((genre, index) => {
                if(!genre.genreCover) return
                return(
                <Card key={index}>
                    <GenreTitle>{genre.genreName}</GenreTitle>
                    <CoverWrapper>
                        <CoverImage src={genre.genreCover} alt={`${genre.genreName} cover`} />
                    </CoverWrapper>
                </Card>
                )})}
        </Grid>
    );
};
