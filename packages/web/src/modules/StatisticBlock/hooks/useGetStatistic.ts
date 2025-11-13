import { getAvarageListens, getListensByDay, getMedianListens, getPeakAndSilentDay } from "@api/analytics";
import { useQuery } from "@tanstack/react-query";

export const useGetListensByDay = () => {
    return useQuery({
        queryKey: ["get-listens"],
        queryFn: () => getListensByDay()
    });
};

export const useGetAvarageListens = () => {
    return useQuery({
        queryKey: ["get-avg-listens"],
        queryFn: () => getAvarageListens()
    });
};

export const useGetMedianListens = () => {
    return useQuery({
        queryKey: ["get-median-listens"],
        queryFn: () => getMedianListens()
    });
};

export const useGetPeakAndSilentDay = () => {
    return useQuery({
        queryKey: ["get-peak-and-silent-day"],
        queryFn: () => getPeakAndSilentDay()
    });
};