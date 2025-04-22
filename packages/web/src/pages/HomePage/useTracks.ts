import { useQuery } from "@tanstack/react-query";
import { getTracks } from "@api/track.api";
import { ISong } from "types/song/song.type";

export const useTracks = () => {
    return useQuery<ISong[]>({
        queryKey: ["tracks"],
        queryFn: getTracks,
    });
};