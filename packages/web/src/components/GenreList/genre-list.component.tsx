import { FC } from "react";
import { GenreResponseType } from "types/song/genreResponseType";
import { Card, CoverImage, CoverWrapper, GenreTitle, Grid } from "./genre-list.component.style";
import { getRandomGradient } from "helpers/get-random-gradient";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";

interface IGenreListProps {
    genres: GenreResponseType[];
}

export const GenreList: FC<IGenreListProps> = ({ genres }) => {
    const navigate = useNavigate()

    return (
        <Grid>
            {genres.map((genre, index) => {
                if (!genre.genreCover) return null;

                const gradient = getRandomGradient();

                return (
                    <Card key={index} background={gradient} onClick={() => navigate(ROUTES.GENRE_SONGS.replace(":genre", genre.genreName))}>
                        <GenreTitle>{genre.genreName}</GenreTitle>
                        <CoverWrapper>
                            <CoverImage src={genre.genreCover} alt={`${genre.genreName} cover`} />
                        </CoverWrapper>
                    </Card>
                );
            })}
        </Grid>
    );
};
