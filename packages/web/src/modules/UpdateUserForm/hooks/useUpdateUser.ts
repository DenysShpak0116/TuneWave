import { useMutation } from "@tanstack/react-query";
import { updateUser } from "@api/user.api"

export const useUpdateUser = () => {
    return useMutation({
        mutationFn: ({ id, profileInfo, username }: { id: string; profileInfo: string; username: string }) =>
            updateUser(id, profileInfo, username),
    });
};
