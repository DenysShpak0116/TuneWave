import { MainLayout } from "@ui/layout/main-layout";
import { FC, useState } from "react";
import { useParams } from "react-router-dom";
import { useGetCollection, useHasCollectionHaveAllVectors } from "./hooks/useGetCollection";
import { Loader } from "@ui/Loader/loader.component";
import { TrackInformationLayout } from "@ui/layout/TrackInformation/track-information-layout";
import { CollectionLogo } from "@components/CollectionLogo/collection-logo.component";
import { getTotalDuration } from "helpers/song/getTotalDuration";
import { TrackDetails } from "@modules/TrackInformation/components/TrackDetails/track-details.component";
import { TrackPagePlayer } from "@components/TrackPagePlayer/track-page-player.component";
import { CollectionSongs } from "@components/CollectionSongs/collection-songs.component";
import { useAuthStore } from "@modules/LoginForm/store/store";
import { useGetCollectionSongs } from "@components/CollectionSongs/hooks/useCollectionSongsю";

export const CollectionPage: FC = () => {
    const userId = useAuthStore(store => store.user?.id)
    const { id } = useParams();
    const [params, setParams] = useState({ search: "", sortBy: "createdAt", order: "desc" });

    const { data: collection, isLoading, refetch } = useGetCollection(id!);
    const { data: collectionSongs, isLoading: loadingSongs } = useGetCollectionSongs(id, params);
    const { data: hasAllVectors = false, isLoading: IsVectorLoading } = useHasCollectionHaveAllVectors(id!)

    const updateSearchSortParams = (updatedParams: { search: string; sortBy: string; order: string }) => {
        setParams(updatedParams);
    };

    if (IsVectorLoading || isLoading || !collection)
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        );

    const isMainUserCollection = userId === collection.user.id;
    const total = getTotalDuration(collection.collectionSongs);

    return (
        <MainLayout>
            <TrackInformationLayout>
                <CollectionLogo
                    hasAllVectors={hasAllVectors}
                    logo={collection.coverUrl}
                    collectionSongs={collection.collectionSongs}
                    collectionId={collection.id}
                />
                <TrackDetails
                    isMainUser={isMainUserCollection}
                    collectionId={collection.id}
                    duration={total}
                    date={collection.createdAt}
                    username={collection.user.username}
                    userId={collection.user.id}
                    type="collection"
                    collectionName={collection.title}
                    collectionDescription={collection.description}
                />
                {collection.collectionSongs.length > 0 && (
                    <TrackPagePlayer song={collection.collectionSongs[0]} />
                )}
                <CollectionSongs
                    refetchFn={updateSearchSortParams}
                    songs={collectionSongs!}
                />
            </TrackInformationLayout>
        </MainLayout>
    );
};

