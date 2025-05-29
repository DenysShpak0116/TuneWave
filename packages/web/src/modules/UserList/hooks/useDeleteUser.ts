import { unfollowUser } from "@api/user.api"
import { useMutation } from "@tanstack/react-query"
import toast from "react-hot-toast"

export const useUnfollow = (onSuccess?: () => void) => {
    return useMutation({
        mutationFn: unfollowUser,
        onSuccess: () => {
            toast.success("Користувача усішно видалено зі списку")
            if (onSuccess) onSuccess()
        }
    })
}