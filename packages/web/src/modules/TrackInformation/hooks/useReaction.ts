import { getUserReaction, reactSong } from "@api/track.api";
import { useMutation, useQuery } from "@tanstack/react-query";
import toast from "react-hot-toast";

export const useReaction = () => {
    return useMutation({
        mutationFn: async ({ songId, reactionType, userId }:
            { songId: string, reactionType: "like" | "dislike", userId: string }) =>
            await reactSong(songId, reactionType, userId),

        onSuccess: () => {
            toast.success("Реація успішно додана")
        },
        onError: (err) => {
            toast.error(`Помилка додавання реакції ${err}`);
        }
    });
};

export const useGetUserReaction = (songId: string, userId: string) => {
    return useQuery<"like" | "dislike", "none">({
        queryKey: ["user-reaction", songId, userId],
        queryFn: () => getUserReaction(songId, userId),
        enabled: !!songId && !!userId,
    });
};