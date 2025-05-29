import { FC, useEffect, useState } from "react";
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
    refetchFn: (params: { search: string; sortBy: string; order: string }) => void;
}

export const CollectionSongs: FC<ICollectionSongsProps> = ({ songs, refetchFn }) => {
    const [activeSongId, setActiveSongId] = useState<string | null>(null);
    const pausePlayer = usePlayerStore((state) => state.pausePlayer);

    const [search, setSearch] = useState('');
    const [sortBy, setSortBy] = useState<'title' | 'added_at'>('added_at');
    const [order, setOrder] = useState<'asc' | 'desc'>('desc');

    useEffect(() => {
        refetchFn({ search, sortBy, order });
    }, [search, sortBy, order]);

    const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        setSearch(e.target.value);
    };

    const handleSortToggle = () => {
        setSortBy(prev => prev === 'title' ? 'added_at' : 'title');
    };

    const handleOrderToggle = () => {
        setOrder(prev => prev === 'asc' ? 'desc' : 'asc');
    };

    const renderPlaylist = () => (
        <>
            <PlaylistHeader>
                <PlaylistTitle>Плейлист</PlaylistTitle>
                <PlaylistActions>
                    <PlayListActionItem>
                        <PlayListActionIcon src={searchIcon} />
                        <input
                            type="text"
                            placeholder="Пошук"
                            value={search}
                            onChange={handleSearchChange}
                            style={{
                                background: "transparent",
                                border: "none",
                                color: "white",
                                outline: "none",
                                fontSize: 14,
                            }}
                        />
                    </PlayListActionItem>
                    <PlayListActionItem onClick={handleSortToggle}>
                        <PlayListActionIcon src={sortIcon} />
                        <p>Сортувати: {sortBy === 'title' ? 'Назва' : 'Дата'}</p>
                    </PlayListActionItem>
                    <PlayListActionItem onClick={handleOrderToggle}>
                        <PlayListActionIcon src={sortIcon} />
                        <p>Порядок: {order === 'asc' ? '↑' : '↓'}</p>
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
                    refetchFn={() => refetchFn({ search, sortBy, order })}
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
            {!songs || songs.length <= 0 ? (
                <>
                    <PlaylistHeader>
                        <PlaylistTitle>Плейлист</PlaylistTitle>
                        <PlaylistActions>
                            <PlayListActionItem>
                                <PlayListActionIcon src={searchIcon} />
                                <input
                                    type="text"
                                    placeholder="Пошук"
                                    value={search}
                                    onChange={handleSearchChange}
                                    style={{
                                        background: "transparent",
                                        border: "none",
                                        color: "white",
                                        outline: "none",
                                        fontSize: 14,
                                    }}
                                />
                            </PlayListActionItem>
                            <PlayListActionItem onClick={handleSortToggle}>
                                <PlayListActionIcon src={sortIcon} />
                                <p>Сортувати: {sortBy === 'title' ? 'Назва' : 'Дата'}</p>
                            </PlayListActionItem>
                            <PlayListActionItem onClick={handleOrderToggle}>
                                <PlayListActionIcon src={sortIcon} />
                                <p>Порядок: {order === 'asc' ? '↑' : '↓'}</p>
                            </PlayListActionItem>
                        </PlaylistActions>
                    </PlaylistHeader>
                    <p>Пісні не знайдено або їх нема</p>
                </>

            ) : (
                renderPlaylist()
            )}
        </CollectionSongsContainer>
    );
};
