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
import { parseDate } from "helpers/date-parse";
import playIcon from "@assets/images/ic_play.png";
import pauseIcon from "@assets/images/ic_pause.png";
import { usePlayerStore } from "@modules/Player/store/player.store";
import { TrackData } from "@components/SongCards/song-cards.component";
import { useParams } from "react-router-dom";
import { useDeleteTrackFromCollection } from "../hooks/useDeleteSongFromCollection";

interface Props {
    collectionSongs: ISong[]
    song: ISong;
    index: number;
    active: boolean;
    onActivate: () => void;
    refetchFn: () => void
}

export const CollectionSongRow: FC<Props> = ({ song, index, active, onActivate, refetchFn, collectionSongs }) => {
    const { id } = useParams()
    const { mutate: removeSongFromCollection } = useDeleteTrackFromCollection()
    const [isHovered, setIsHovered] = useState(false);
    const [showOptions, setShowOptions] = useState(false);
    const optionsRef = useRef<HTMLDivElement>(null);
    const { setPlaylist, setTrack, setIsPlaying } = usePlayerStore();

    const handlePlay = (trackData: TrackData) => {
        onActivate();
        const mappedPlaylist = collectionSongs.map(song => ({
            id: song.id,
            file: encodeURI(song.songUrl),
            logo: song.coverUrl,
            title: song.title,
            artist: song.user.username,
        }));

        setPlaylist(mappedPlaylist);
        setTrack(trackData);
        setIsPlaying(true);
    };
    
    const handleDelete = () => {
        removeSongFromCollection({ songId: song.id, collectionId: id! })
        refetchFn()
        setShowOptions(false);
    };

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (
                optionsRef.current &&
                !optionsRef.current.contains(event.target as Node)
            ) {
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
            active={active}
            isHovered={isHovered}
            onMouseEnter={() => setIsHovered(true)}
            onMouseLeave={() => setIsHovered(false)}
        >
            <IndexBox>
                {(isHovered || active) ? (
                    <PlayListActionIcon
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
                        src={active ? pauseIcon : playIcon}
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
            <Duration>{song.duration}</Duration>
            <Options onClick={(e) => {
                console.log("dadas");

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
