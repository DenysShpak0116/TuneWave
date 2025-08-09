import { updateTrack } from "@api/track.api";
import { useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";

export const useUpdateTrack = () => {
    return useMutation({
        mutationFn: ({ songId, formData }: { songId: string; formData: FormData }) =>
            updateTrack(songId, formData),
        onSuccess: () => {
            toast.success("Трек успішно оновлено")
        },
        onError: (err) => {
            toast.error(`Помилка при оновленні треку ${err}`);
        },
    });
};