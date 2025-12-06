import { getPopularTracks, getTracksWithPopular } from "@api/analytics";
import { useQuery } from "@tanstack/react-query";
import { ITopTrackStatistic } from "types/song/song-statistic";

export const useGetTopTracks = () => {
    return useQuery<ITopTrackStatistic[]>({
        queryKey: ["get-popular-tracks"],
        queryFn: () => getPopularTracks()
    });
};

export const useGetTracksWithPopular = () => {
    return useQuery({
        queryKey: ["get-tracks-with-popular"],
        queryFn: () => getTracksWithPopular()
    });
};