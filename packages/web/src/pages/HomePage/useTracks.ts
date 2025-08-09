import { useQuery } from "@tanstack/react-query";
import { getTracks } from "@api/track.api";
import { ISong } from "types/song/song.type";
import { GetTracksParams } from "types/song/getTrackParams";

export const useTracks = (params?: GetTracksParams) => {
    return useQuery<ISong[]>({
        queryKey: ["tracks", params ?? {}],
        queryFn: () => getTracks(params ?? {}),
    });
};