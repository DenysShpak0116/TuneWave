import { getUserCollections } from "@api/user.api";
import { useQuery } from "@tanstack/react-query";
import { ICollection } from "types/collections/collection.type";

export const useGetUserCollections = (id: string) => {
    return useQuery<ICollection[]>({
        queryKey: ["user-collections", id],
        queryFn: () => getUserCollections(id),
        enabled: !!id,
    });
};