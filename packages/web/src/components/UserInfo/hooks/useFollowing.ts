import { followUser, isFollowed } from "@api/user.api";
import { useMutation, useQuery } from "@tanstack/react-query";
import toast from "react-hot-toast";

export const useIsFollowed = (userId: string) => {
    return useQuery({
        queryKey: ["is-followed", userId],
        queryFn: () => isFollowed(userId),
        enabled: !!userId,
    });
};

export const useFollow = (onSuccess?: () => void) => {
    return useMutation({
        mutationFn: followUser,
        onSuccess: () => {
            toast.success("Ви підписалися на користувача")
            if (onSuccess) onSuccess()
        }
    });
};