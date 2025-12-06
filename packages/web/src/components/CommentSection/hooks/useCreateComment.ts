import { createComment } from "@api/comment.api"
import { useMutation } from "@tanstack/react-query"

export const useCreateComment = () => {
    return useMutation({
        mutationFn: ({
            userId,
            songId,
            content,
        }: {
            userId: string
            songId: string
            content: string
        }) => createComment(userId, songId, content),
    })
}