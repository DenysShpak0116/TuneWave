import { FC } from "react";
import { useNavigate } from "react-router-dom";
import { ISong } from "types/song/song.type";
import { ICollection } from "types/collections/collection.type";
import { usePlayerStore } from "@modules/Player/store/player.store";
import { ROUTES } from "pages/router/consts/routes.const";
import playIcon from "@assets/images/ic_play.png";

import {
    Wrapper,
    SongCardsContainer,
    SongCard,
    SongImage,
    SongTitle,
    SongArtist,
    SongsText,
    ImageWrapper,
    PlayIcon,
} from "./song-cards.style";

interface ISongCardsProps {
    songs?: ISong[];
    collections?: ICollection[];
    text: string;
}

export interface TrackData {
    trackId: string;
    trackUrl: string;
    trackLogo: string;
    trackName: string;
    trackArtist: string;
}

export const SongCards: FC<ISongCardsProps> = ({ songs, collections, text }) => {
    const navigate = useNavigate();
    const setTrack = usePlayerStore((state) => state.setTrack);

    const handlePlay = (trackData: TrackData) => {
        setTrack(trackData);
    };

    const handleNavigate = (path: string) => {
        navigate(path);
    };

    const renderSongCard = (song: ISong) => (
        <SongCard
            key={song.id}
            onClick={() => handleNavigate(ROUTES.TRACK_PAGE.replace(":id", song.id))}
        >
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
    );

    const renderCollectionCard = (collection: ICollection) => (
        <SongCard
            key={collection.id}
            onClick={() => handleNavigate(ROUTES.COLLECTION_PAGE.replace(":id", collection.id))}
        >
            <ImageWrapper>
                <SongImage src={collection.coverUrl} alt={collection.title} />
            </ImageWrapper>
            <SongTitle>{collection.title}</SongTitle>
        </SongCard>
    );

    const renderCards = () => {
        if (songs) return songs.map(renderSongCard);
        if (collections) return collections.map(renderCollectionCard);
        return <p>Loading...</p>;
    };

    return (
        <Wrapper>
            <SongsText>{text}</SongsText>
            <SongCardsContainer>{renderCards()}</SongCardsContainer>
        </Wrapper>
    );
};
