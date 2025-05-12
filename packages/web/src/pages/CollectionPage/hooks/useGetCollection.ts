import { getCollectionById } from "@api/collection.api"
import { useQuery } from "@tanstack/react-query"
import { ICollection } from "types/collections/collection.type"

export const useGetCollection = (id: string) => {
    return useQuery<ICollection>({
        queryKey: ["collection", id],
        queryFn: () => getCollectionById(id),
        enabled: !!id
    })
}