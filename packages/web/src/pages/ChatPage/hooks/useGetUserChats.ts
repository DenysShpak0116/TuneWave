import { getChatPreviews } from "@api/user.api"
import { useQuery } from "@tanstack/react-query"
import { IChatPreviewType } from "types/chat/chat-preview"

export const useGetUserChats = () => {
    return useQuery<IChatPreviewType[]>({
        queryKey: ["chat-previews"],
        queryFn: () => getChatPreviews(),
    })
}