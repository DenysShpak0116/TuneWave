import { getCollectionSongs } from "@api/collection.api";
import { useQuery } from "@tanstack/react-query";
import { CollectionSongsParams } from "types/collections/collectionSongsParams.type";
import { ISong } from "types/song/song.type";

export const useGetCollectionSongs = (collectonId: string, params?: CollectionSongsParams) => {
    return useQuery<ISong[]>({
        queryKey: ["collection-songs", collectonId, params ?? {}],
        queryFn: () => getCollectionSongs(collectonId, params ?? {}),
    });
};