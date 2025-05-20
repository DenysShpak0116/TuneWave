import { FC, useState } from "react";
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
} from "../collection-songs.style";
import { parseDate } from "helpers/date-parse";
import playIcon from "@assets/images/ic_play.png"
import pauseIcon from "@assets/images/ic_pause.png"
import { usePlayerStore } from "@modules/Player/store/player.store";
import { TrackData } from "@components/SongCards/song-cards.component";

interface Props {
    song: ISong;
    index: number;
    active: boolean;
    onActivate: () => void;
}

export const CollectionSongRow: FC<Props> = ({ song, index, active, onActivate }) => {
    const [isHovered, setIsHovered] = useState(false);
    const { setTrack, setIsPlaying } = usePlayerStore();

    const handlePlay = (trackData: TrackData) => {
        if (active) {
            setIsPlaying(false);
        }
        onActivate()
        setTrack(trackData);
    };

    return (
        <TableRow
            active={active}
            isHovered={isHovered}
            onMouseEnter={() => setIsHovered(true)}
            onMouseLeave={() => setIsHovered(false)}
        >
            <IndexBox>
                {(isHovered || active) ? <PlayListActionIcon
                    onClick={(e) => {
                        e.stopPropagation()
                        handlePlay({
                            trackId: song.id,
                            trackUrl: encodeURI(song.songUrl),
                            trackLogo: song.coverUrl,
                            trackName: song.title,
                            trackArtist: song.user.username,
                        })
                    }}
                    src={active ? pauseIcon : playIcon} />
                    : index + 1}
            </IndexBox>
            <CoverAndInfo>
                <Cover src={song.coverUrl} alt={song.title} />
                <SongTextInfo>
                    <Title>{song.title}</Title>
                </SongTextInfo>
            </CoverAndInfo>
            <DateAdded>{parseDate(song.createdAt)}</DateAdded>
            <Duration>{song.duration}</Duration>
            <Options>⋯</Options>
        </TableRow>
    );
};
