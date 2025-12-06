import { getCommonCombos, getRareCombos } from "@api/analytics";
import { useQuery } from "@tanstack/react-query";

export const useGetCommonCombinations = () => {
    return useQuery({
        queryKey: ["get-common-combinations"],
        queryFn: () => getCommonCombos()
    });
};

export const useGetRareCombinations = () => {
    return useQuery({
        queryKey: ["get-rare-combinations"],
        queryFn: () => getRareCombos()
    });
};

