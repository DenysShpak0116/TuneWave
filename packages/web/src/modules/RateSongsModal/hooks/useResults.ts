import { createResult } from '@api/result.api';
import { useMutation } from '@tanstack/react-query';
import { IResultType } from 'types/results/result.type';

export const useCreateResult = () => {
    return useMutation({
        mutationFn: ({ collectionId, results }: { collectionId: string; results: IResultType[] }) =>
            createResult(collectionId, results),
    });
};