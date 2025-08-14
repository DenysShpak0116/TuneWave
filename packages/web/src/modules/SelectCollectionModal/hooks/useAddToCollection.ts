import { addSongToCollection } from "@api/track.api";
import { useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";

export const useAddSongToCollection = () => {
    return useMutation({
        mutationFn: async ({ trackId, collectionId }: { trackId: string; collectionId: string }) =>
            await addSongToCollection(trackId, collectionId),
        onSuccess: () => {
            toast.success("Пісня успішно додана")
        },
        onError: (err) => {
            toast.error(`Помилка додавання пісні ${err}`);
        }
    });
}