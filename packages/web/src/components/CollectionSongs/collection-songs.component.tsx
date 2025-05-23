import { FC, useState } from "react";
import { ISong } from "types/song/song.type";
import {
    CollectionSongsContainer,
    TableHeader,
    PlaylistHeader,
    PlaylistTitle,
    PlaylistActions,
    PlayListActionIcon,
    PlayListActionItem,
} from "./collection-songs.style";
import searchIcon from "@assets/images/ic_search.png";
import sortIcon from "@assets/images/ic_sort.png";
import { CollectionSongRow } from "./components/collection-song-row.component";
import { usePlayerStore } from "@modules/Player/store/player.store";

interface ICollectionSongsProps {
    songs: ISong[];
    refetchFn: () => void;
}

export const CollectionSongs: FC<ICollectionSongsProps> = ({ songs, refetchFn }) => {
    const [activeSongId, setActiveSongId] = useState<string | null>(null);
    const pausePlayer = usePlayerStore((state) => state.pausePlayer);

    const renderPlaylist = () => (
        <>
            <PlaylistHeader>
                <PlaylistTitle>Плейлист</PlaylistTitle>
                <PlaylistActions>
                    <PlayListActionItem>
                        <PlayListActionIcon src={searchIcon} />
                        <p>Пошук</p>
                    </PlayListActionItem>
                    <PlayListActionItem onClick={pausePlayer}>
                        <PlayListActionIcon src={searchIcon} />
                        <p>Пауза</p>
                    </PlayListActionItem>
                    <PlayListActionItem>
                        <PlayListActionIcon src={sortIcon} />
                        <p>Сортувати</p>
                    </PlayListActionItem>
                </PlaylistActions>
            </PlaylistHeader>

            <TableHeader>
                <div>#</div>
                <div>Назва</div>
                <div>Дата додавання</div>
                <div>Час</div>
                <div style={{ textAlign: "center" }}>Дії</div>
            </TableHeader>

            {songs.map((song, index) => (
                <CollectionSongRow
                    collectionSongs={songs}
                    refetchFn={refetchFn}
                    key={song.id}
                    song={song}
                    index={index}
                    active={song.id === activeSongId}
                    onActivate={() => setActiveSongId(song.id)}
                />
            ))}
        </>
    );

    return (
        <CollectionSongsContainer>
            {songs.length > 0 ? (
                renderPlaylist()
            ) : (
                <p style={{ textAlign: "center", }}>У цьому плейлисті поки немає пісень.</p>
            )}
        </CollectionSongsContainer>
    );
};
