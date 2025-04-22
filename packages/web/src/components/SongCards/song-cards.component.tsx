import { FC } from "react";
import { ISong } from "types/song/song.type";
import {
    Wrapper,
    SongCardsContainer,
    SongCard,
    SongImage,
    SongTitle,
    SongArtist,
    SongsText
} from "./song-cards.style";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";

interface ISongCardsProps {
    songs: ISong[];
}

export const SongCards: FC<ISongCardsProps> = ({ songs }) => {
    const navigate = useNavigate()

    const handleSongCardClick = (id: string) => {
        navigate(ROUTES.TRACK_PAGE.replace(":id", id))
    }

    return (
        <Wrapper>
            <SongsText>ПОПУЛЯРНІ МУЗИЧНІ ТВОРИ</SongsText>
            <SongCardsContainer>
                {songs.map((song) => (
                    <SongCard key={song.id} onClick={() => handleSongCardClick(song.id)}>
                        <SongImage src={song.coverUrl} alt={song.title} />
                        <SongTitle>{song.title}</SongTitle>
                        <SongArtist>{song.user.username}</SongArtist>
                    </SongCard>
                ))}
            </SongCardsContainer>
        </Wrapper>
    );
};
