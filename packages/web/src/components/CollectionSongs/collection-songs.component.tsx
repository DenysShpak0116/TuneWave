import { FC } from "react";
import { ISong } from "types/song/song.type";
import {
    CollectionSongsContainer,
    TableHeader,
    TableRow,
    IndexBox,
    CoverAndInfo,
    Cover,
    SongTextInfo,
    Title,
    Album,
    DateAdded,
    Duration,
    Options
} from "./collection-songs.style";
import { parseDate } from "helpers/date-parse";

interface ICollectionSongsProps {
    songs: ISong[];
    activeSongId?: string;
}

export const CollectionSongs: FC<ICollectionSongsProps> = ({ songs, activeSongId }) => {
    return (
        <CollectionSongsContainer>
            <TableHeader>
                <div>#</div>
                <div>Назва</div>
                <div>Альбом</div>
                <div>Дата додавання</div>
                <div>Час</div>
                <div></div>
            </TableHeader>
            {songs.map((song, index) => (
                <TableRow key={song.id} active={song.id === activeSongId}>
                    <IndexBox>{index + 1}</IndexBox>
                    <CoverAndInfo>
                        <Cover src={song.coverUrl} alt={song.title} />
                        <SongTextInfo>
                            <Title>{song.title}</Title>
                            {/* <Author>{song.authors.map(a => a.name).join(", ")}</Author> */}
                        </SongTextInfo>
                    </CoverAndInfo>
                    <Album>123</Album>
                    <DateAdded>{parseDate(song.createdAt)}</DateAdded>
                    <Duration>{song.duration}</Duration>
                    <Options>⋯</Options>
                </TableRow>
            ))}
        </CollectionSongsContainer>
    );
};
