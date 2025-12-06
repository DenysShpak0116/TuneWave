import { getCollections } from "@api/collection.api";
import { useQuery } from "@tanstack/react-query";
import { ICollection } from "types/collections/collection.type";
import { GetCollectionParams } from "types/collections/getCollectionParams";

export const useGetTopCollections = (params?: GetCollectionParams) => {
    return useQuery<ICollection[]>({
        queryKey: ["collections", params ?? {}],
        queryFn: () => getCollections(params ?? {}),
    });
};