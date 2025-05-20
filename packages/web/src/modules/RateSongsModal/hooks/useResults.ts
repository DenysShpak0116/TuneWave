import { createResult, getUserResults } from '@api/result.api';
import { useMutation, useQuery } from '@tanstack/react-query';
import { IResultType } from 'types/results/result.type';

export const useCreateResult = () => {
    return useMutation({
        mutationFn: ({ collectionId, results }: { collectionId: string; results: IResultType[] }) =>
            createResult(collectionId, results),
    });
};

export const useGetUserResults = (collectionId: string) => {
    return useQuery({
        queryKey: ['userResults', collectionId],
        queryFn: () => getUserResults(collectionId),
        enabled: !!collectionId,
    });
};