import { getGenres } from "@api/track.api"
import { useQuery } from "@tanstack/react-query"
import { GenreResponseType } from "types/song/genreResponseType"

export const useGetGenres = () => {
    return useQuery<GenreResponseType[]>({
        queryKey: ["genres"],
        queryFn: () => getGenres(),
    })
}