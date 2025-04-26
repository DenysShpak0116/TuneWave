import { $api } from "./base.api"

export const createComment = async (userId: string, songId: string, content: string) => {
    const { data } = await $api.post("/comments", { content, userId, songId })
    return data
}

export const deleteComment = async (id: string) => {
    return await $api.delete(`/comments/${id}`)
}