import { useMutation, useQuery } from "@tanstack/react-query";
import { deleteTrack, getTrackById } from "@api/track.api";
import { ISong } from "types/song/song.type";
import toast from "react-hot-toast";

export const useGetTrack = (id: string) => {
    return useQuery<ISong>({
        queryKey: ["track", id],
        queryFn: () => getTrackById(id),
        enabled: !!id,
    })
}

export const useDeleteTrack = () => {
    return useMutation({
        mutationFn: deleteTrack,
        onSuccess: () => {
            toast.success("Пісня видалена успішно")
        }
    })
}