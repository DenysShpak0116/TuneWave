import { useMutation } from "@tanstack/react-query";
import { deleteComment } from "@api/comment.api";

export const useDeleteComment = () => {
    return useMutation({
        mutationFn: (id: string) => deleteComment(id),
    });
};