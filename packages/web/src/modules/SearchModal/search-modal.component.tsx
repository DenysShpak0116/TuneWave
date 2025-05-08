import { FC, useState } from "react";
import { useDebounce } from "use-debounce";
import {
    ModalContent,
    Overlay,
    Header,
    SearchInput,
    FilterButton,
    ScrollContainer,
    SongList,
    FilterContainer,
    FilterLabel,
    SelectInput
} from "./search-modal.style";
import { ISong } from "types/song/song.type";
import { useTracks } from "pages/HomePage/useTracks";
import { SearchSongCard } from "@components/SearchSongCard/search-song-card.component";

interface ISearchModalProps {
    active: boolean;
    setActive: (value: boolean) => void;
}

type sortByType = 'created_at' | 'title' | 'artist' | 'genre'
type orderType = 'asc' | 'desc'

export const SearchModal: FC<ISearchModalProps> = ({ active, setActive }) => {
    const [query, setQuery] = useState<string>("");
    const [debouncedQuery] = useDebounce(query, 500);
    const [showFilters, setShowFilters] = useState<boolean>(false);
    const [limit, setLimit] = useState<number>(10);
    const [sortBy, setSortBy] = useState<sortByType>('created_at');
    const [order, setOrder] = useState<orderType>('desc');

    const { data: tracks, isLoading } = useTracks({
        search: debouncedQuery,
        sortBy,
        order,
        limit
    });

    return (
        <Overlay $active={active} onClick={() => setActive(false)}>
            <ModalContent $active={active} onClick={(e) => e.stopPropagation()}>
                <Header>
                    <SearchInput
                        placeholder="Шукати пісню"
                        value={query}
                        onChange={(e) => setQuery(e.target.value)}
                    />
                    <FilterButton onClick={() => setShowFilters(!showFilters)}>Фільтри</FilterButton>
                </Header>
                {isLoading && <p>Loading...</p>}
                {showFilters && (
                    <FilterContainer>
                        <div>
                            <FilterLabel>
                                <p>Сортувати</p>
                                <SelectInput value={sortBy} onChange={(e) => setSortBy(e.target.value as sortByType)}>
                                    <option value="created_at">Дата</option>
                                    <option value="title">Назва</option>
                                    <option value="artist">Артист</option>
                                    <option value="genre">Жанр</option>
                                </SelectInput>
                            </FilterLabel>
                        </div>
                        <div>
                            <FilterLabel>
                                <p>Порядок</p>
                                <SelectInput value={order} onChange={(e) => setOrder(e.target.value as orderType)}>
                                    <option value="desc">За спаданням</option>
                                    <option value="asc">За зростанням</option>
                                </SelectInput>
                            </FilterLabel>
                        </div>
                        <div>
                            <FilterLabel>
                                <p>Ліміт</p>
                                <SelectInput value={limit} onChange={(e) => setLimit(e.target.value as any)}>
                                    <option value={1}>1</option>
                                    <option value={10}>10</option>
                                    <option value={20}>20</option>
                                    <option value={40}>40</option>
                                </SelectInput>
                            </FilterLabel>
                        </div>
                    </FilterContainer>
                )}

                {!isLoading && debouncedQuery && tracks?.length === 0 && <p>За запитом "{debouncedQuery}" нічого не знадено</p>}
                {!isLoading && debouncedQuery && tracks && tracks?.length > 0 && (
                    <ScrollContainer>
                        <SongList>
                            {tracks.map((track: ISong) => (
                                <li key={track.id}>
                                    <SearchSongCard song={track} />
                                </li>
                            ))}
                        </SongList>
                    </ScrollContainer>
                )}
            </ModalContent>
        </Overlay>
    );
};
