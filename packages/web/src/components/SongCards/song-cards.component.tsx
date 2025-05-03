import { FC } from "react";
import { ISong } from "types/song/song.type";
import {
    Wrapper,
    SongCardsContainer,
    SongCard,
    SongImage,
    SongTitle,
    SongArtist,
    SongsText,
    ImageWrapper,
    PlayIcon
} from "./song-cards.style";
import { useNavigate } from "react-router-dom";
import { ROUTES } from "pages/router/consts/routes.const";
import playIcon from "@assets/images/ic_play.png"
import { usePlayerStore } from "@modules/Player/store/player.store";

interface ISongCardsProps {
    songs: ISong[];
}

export interface TrackData {
    trackId: string;
    trackUrl: string;
    trackLogo: string;
    trackName: string;
    trackArtist: string;
}

export const SongCards: FC<ISongCardsProps> = ({ songs }) => {
    const setTrack = usePlayerStore((state) => state.setTrack);

    const handlePlay = (trackData: TrackData) => {
        setTrack(trackData);
    };

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
                        <ImageWrapper>
                            <SongImage src={song.coverUrl} alt={song.title} />
                            <PlayIcon
                                className="play-icon"
                                src={playIcon}
                                alt="Play"
                                onClick={(e) => {
                                    e.stopPropagation();
                                    handlePlay({
                                        trackId: song.id,
                                        trackUrl: encodeURI(song.songUrl),
                                        trackLogo: song.coverUrl,
                                        trackName: song.title,
                                        trackArtist: song.user.username,
                                    });
                                }}
                            />
                        </ImageWrapper>
                        <SongTitle>{song.title}</SongTitle>
                        <SongArtist>{song.user.username}</SongArtist>
                    </SongCard>
                ))}
            </SongCardsContainer>
        </Wrapper>
    );
};
