import { useMutation, useQuery } from "@tanstack/react-query";
import { deleteTrack, getTrackById, getTrackCommentsById } from "@api/track.api";
import { ISong } from "types/song/song.type";
import toast from "react-hot-toast";
import { IComment } from "types/comments/comment.type";

export const useGetTrack = (id: string) => {
    return useQuery<ISong>({
        queryKey: ["track", id],
        queryFn: () => getTrackById(id),
        enabled: !!id,
    })
}

export const useGetTrackComments = (id: string) => {
    return useQuery<IComment[]>({
        queryKey: ["track-comments", id],
        queryFn: () => getTrackCommentsById(id),
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