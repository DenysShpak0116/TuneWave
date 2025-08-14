import { FC, useState, useRef, useEffect } from "react";
import { ISong } from "types/song/song.type";
import {
    TableRow,
    IndexBox,
    CoverAndInfo,
    Cover,
    SongTextInfo,
    Title,
    DateAdded,
    Duration,
    Options,
    PlayListActionIcon,
    OptionsPopup,
    OptionButton,
} from "../collection-songs.style";
import { parseDate, parseTime } from "helpers/date-parse";
import playIcon from "@assets/images/ic_play.png";
import pauseIcon from "@assets/images/ic_pause.png";
import { usePlayerStore } from "@modules/Player/store/player.store";
import { TrackData } from "@components/SongCards/song-cards.component";
import { useParams } from "react-router-dom";
import { useDeleteTrackFromCollection } from "../hooks/useDeleteSongFromCollection";

interface Props {
    collectionSongs: ISong[];
    song: ISong;
    index: number;
    refetchFn: () => void;
}

export const CollectionSongRow: FC<Props> = ({ song, index, refetchFn, collectionSongs }) => {
    const { id } = useParams();
    const { mutate: removeSongFromCollection } = useDeleteTrackFromCollection();
    const [isHovered, setIsHovered] = useState(false);
    const [showOptions, setShowOptions] = useState(false);
    const optionsRef = useRef<HTMLDivElement>(null);
    const { trackId: currentTrackId, isPlaying, pausePlayer, setTrack, setPlaylist } = usePlayerStore();
    const isActive = song.id === currentTrackId;

    const handlePlay = () => {
        if (isActive && isPlaying) {
            pausePlayer();
        } else {
            const trackData: TrackData = {
                trackId: song.id,
                trackUrl: encodeURI(song.songUrl),
                trackLogo: song.coverUrl,
                trackName: song.title,
                trackArtist: song.user.username,
            };
            const playlist = collectionSongs.map(s => ({
                id: s.id,
                file: encodeURI(s.songUrl),
                logo: s.coverUrl,
                title: s.title,
                artist: s.user.username,
            }));
            setPlaylist(playlist);
            setTrack(trackData);
        }
    };

    const handleDelete = () => {
        removeSongFromCollection({ songId: song.id, collectionId: id! });
        refetchFn();
        setShowOptions(false);
    };

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (optionsRef.current && !optionsRef.current.contains(event.target as Node)) {
                setShowOptions(false);
            }
        };

        document.addEventListener("mousedown", handleClickOutside);
        return () => {
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, []);

    return (
        <TableRow
            active={isActive}
            isHovered={isHovered}
            onMouseEnter={() => setIsHovered(true)}
            onMouseLeave={() => setIsHovered(false)}
        >
            <IndexBox>
                {(isHovered || isActive) ? (
                    <PlayListActionIcon
                        onClick={handlePlay}
                        src={isActive && isPlaying ? pauseIcon : playIcon}
                    />
                ) : (
                    index + 1
                )}
            </IndexBox>
            <CoverAndInfo>
                <Cover src={song.coverUrl} alt={song.title} />
                <SongTextInfo>
                    <Title>{song.title}</Title>
                </SongTextInfo>
            </CoverAndInfo>
            <DateAdded>{parseDate(song.createdAt)}</DateAdded>
            <Duration>{parseTime(song.duration)}</Duration>
            <Options onClick={(e) => {
                e.stopPropagation();
                setShowOptions(!showOptions);
            }}>
                ⋯
                {showOptions && (
                    <OptionsPopup ref={optionsRef}>
                        <OptionButton onClick={handleDelete}>Видалити</OptionButton>
                    </OptionsPopup>
                )}
            </Options>
        </TableRow>
    );
};