import { addCollectionToUser } from "@api/collection.api";
import { useMutation } from "@tanstack/react-query";
import toast from "react-hot-toast";

export const useAddCollectionToUser = () => {
    return useMutation({
        mutationFn: (collectionId: string) => addCollectionToUser(collectionId),
        onSuccess: () => {
            toast.success("Колекція додана успішно")
        }
    });
};