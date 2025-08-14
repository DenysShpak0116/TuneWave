import { useMutation } from "@tanstack/react-query";
import { updateUser, updateUserAvatar } from "@api/user.api"
import toast from "react-hot-toast";

export const useUpdateUser = () => {
    return useMutation({
        mutationFn: async ({ id, profileInfo, username }: { id: string; profileInfo: string; username: string }) =>
            await updateUser(id, profileInfo, username),
    });
};

export const useUpdateAvatar = () =>{
    return useMutation({
        mutationFn: async (formData: FormData) => {
            await updateUserAvatar(formData);
        },
        onSuccess: () =>{
            toast.success("Аватар користувача успішно змінено")
        }
    })
}