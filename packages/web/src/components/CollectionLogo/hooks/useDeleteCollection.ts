import { deleteCollection } from "@api/collection.api"
import { useMutation } from "@tanstack/react-query"
import toast from "react-hot-toast"

export const useDeleteCollection = () => {
    return useMutation({
        mutationFn: deleteCollection,
        onSuccess: () => {
            toast.success("Колекція видалена успішно")
        }
    })
}