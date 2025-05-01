import { useMutation } from "@tanstack/react-query";
import { updateUser } from "@api/user.api"

export const useUpdateUser = () => {
    return useMutation({
        mutationFn: async ({ id, profileInfo, username }: { id: string; profileInfo: string; username: string }) =>
            await updateUser(id, profileInfo, username),
    });
};
