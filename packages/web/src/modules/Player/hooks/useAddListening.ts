import { addListening } from "@api/track.api";
import { useMutation } from "@tanstack/react-query";

export const useAddListening = () => {
    return useMutation({
        mutationFn: async ({ songId, userId }: { songId: string, userId: string }) => {
            await addListening(songId, userId)
        },
    });
};