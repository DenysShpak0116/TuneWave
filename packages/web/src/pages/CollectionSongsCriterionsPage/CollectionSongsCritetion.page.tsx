import { FC } from "react";
import { useParams } from "react-router-dom";
import { MainLayout } from "@ui/layout/main-layout";
import { Loader } from "@ui/Loader/loader.component";
import { useGetCollection } from "pages/CollectionPage/hooks/useGetCollection";
import { ISong } from "types/song/song.type";
import { useDeleteVector, useGetVectorBySongId } from "@modules/AddCriterionToSongModal/hooks/useAddVector";
import { DeleteButton, Table, Td, Th } from "./CollectionSongsCritetion.style";

export const CollectionSongsPage: FC = () => {
    const { id: collectionId } = useParams();
    const { data: collection, isLoading } = useGetCollection(collectionId!);
    const deleteVectorMutation = useDeleteVector();

    if (isLoading || !collection) {
        return (
            <MainLayout>
                <Loader />
            </MainLayout>
        );
    }

    const vectorData = collection.collectionSongs.map((song: ISong) => {
        // eslint-disable-next-line react-hooks/rules-of-hooks
        const { data } = useGetVectorBySongId(collectionId!, song.id);
        return {
            songId: song.id,
            songTitle: song.title,
            vectors: data ?? [],
        };
    });

    const uniqueCriteria = Array.from(
        new Set(
            vectorData.flatMap(entry => entry.vectors.map(v => v.criterion))
        )
    );

    const handleDelete = (songId: string) => {
        if (collectionId) {
            deleteVectorMutation.mutate({ collectionId, songId });
        }
    };

    return (
        <MainLayout>
            <Table>
                <thead>
                    <tr>
                        <Th>Пісня</Th>
                        {uniqueCriteria.map((criterion) => (
                            <Th key={criterion}>{criterion}</Th>
                        ))}
                        <Th>Дія</Th>
                    </tr>
                </thead>
                <tbody>
                    {vectorData.map(({ songId, songTitle, vectors }) => (
                        <tr key={songId}>
                            <Td>{songTitle}</Td>
                            {uniqueCriteria.map((criterion) => {
                                const found = vectors.find(v => v.criterion === criterion);
                                return <Td key={criterion}>{found?.mark ?? '-'}</Td>;
                            })}
                            <Td>
                                <DeleteButton onClick={() => handleDelete(songId)}>
                                    Видалити вектори
                                </DeleteButton>
                            </Td>
                        </tr>
                    ))}
                </tbody>
            </Table>
        </MainLayout>
    );
};