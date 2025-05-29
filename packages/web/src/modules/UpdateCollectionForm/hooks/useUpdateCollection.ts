import { updateCollection } from "@api/collection.api";
import { useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";

export const useUpdateTrack = () => {
    return useMutation({
        mutationFn: ({ collectionId, formData }: { collectionId: string; formData: FormData }) =>
            updateCollection(collectionId, formData),
        onSuccess: () => {
            toast.success("Колекцію успішно оновлено")
        },
        onError: (err) => {
            toast.error(`Помилка при оновленні колекції ${err}`);
        },
    });
};