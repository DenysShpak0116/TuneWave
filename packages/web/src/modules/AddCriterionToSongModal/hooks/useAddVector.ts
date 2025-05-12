import { addVector, getVectorBySongId, updateVector } from "@api/vector.api";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { IVectorType } from "types/vectors/vector.type";
interface AddVectorParams {
    collectionId: string;
    songId: string;
    vectors: IVectorType[];
}

export const useAddVector = () => {
    return useMutation({
        mutationKey: ["addVector"],
        mutationFn: ({ collectionId, songId, vectors }: AddVectorParams) =>
            addVector(collectionId, songId, vectors),
    });
};

export const useGetVectorBySongId = (collectionId: string, songId: string) => {
    return useQuery<IVectorType[]>({
        queryKey: ['vectors', collectionId, songId],
        queryFn: () => getVectorBySongId(collectionId, songId),
        enabled: !!collectionId && !!songId,
    })
}

export const useUpdateVector = () => {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: ({
            collectionId,
            songId,
            vectors,
        }: {
            collectionId: string
            songId: string
            vectors: IVectorType[]
        }) => updateVector(collectionId, songId, vectors),

        onSuccess: (_, { collectionId, songId }) => {
            queryClient.invalidateQueries({
                queryKey: ['vectors', collectionId, songId],
            })
        },
    })
}