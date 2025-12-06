import { getByUserId } from "@api/collection.api";
import { useQuery } from "@tanstack/react-query";


export const useUserCollections = (userId: string) => {
    return useQuery({
        queryKey: ["user-collections", userId],
        queryFn: () => getByUserId(),
        enabled: !!userId,
    });
};