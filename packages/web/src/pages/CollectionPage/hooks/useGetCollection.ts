import { getCollectionById } from "@api/collection.api"
import { hasAllVectors } from "@api/vector.api"
import { useQuery } from "@tanstack/react-query"
import { ICollection } from "types/collections/collection.type"

export const useGetCollection = (id: string) => {
    return useQuery<ICollection>({
        queryKey: ["collection", id],
        queryFn: () => getCollectionById(id),
        enabled: !!id
    })
}

export const useHasCollectionHaveAllVectors = (collectionId: string) => {
    return useQuery<boolean>({
        queryKey: ["collection-vectors", collectionId],
        queryFn: () => hasAllVectors(collectionId),
        enabled: !!collectionId
    })
}