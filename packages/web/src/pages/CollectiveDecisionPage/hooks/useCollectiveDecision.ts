import { getCollectiveDecision } from "@api/result.api";
import { useQuery } from "@tanstack/react-query";
import { CollectiveDecisionResponse } from "types/results/collective-decision.type";

export const useCollectiveDecision = (collectionId: string) => {
    return useQuery<CollectiveDecisionResponse>({
        queryKey: ['collective-decision', collectionId],
        queryFn: () => getCollectiveDecision(collectionId),
        enabled: !!collectionId,
    });
};