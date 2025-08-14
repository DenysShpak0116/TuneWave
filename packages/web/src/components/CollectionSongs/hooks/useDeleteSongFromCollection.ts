import { removeSongFromCollection } from "@api/track.api"
import { useMutation } from "@tanstack/react-query"
import toast from "react-hot-toast"

export const useDeleteTrackFromCollection = () => {
    return useMutation({
        mutationFn: ({ songId, collectionId }: { songId: string, collectionId: string }) =>
            removeSongFromCollection(songId, collectionId),
        onSuccess: () => {
            toast.success("Пісня видалена успішно з колекції")
        }
    })
}