import { useQuery } from "@tanstack/react-query";
import { getTrackById } from "@api/track.api";
import { ISong } from "types/song/song.type";

export const useGetTrack = (id: string) => {
    return useQuery<ISong>({
        queryKey: ["track", id],
        queryFn: () => getTrackById(id),
        enabled: !!id,
    })
}