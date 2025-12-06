import { getPublicKey } from "@api/auth.api";
import { useQuery } from "@tanstack/react-query";

export const useGetPublicKey = () => {
    return useQuery({
        queryKey: ["public-key"],
        queryFn: () => getPublicKey(),
    });
};