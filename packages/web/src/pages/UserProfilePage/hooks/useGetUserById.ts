import { getUserById } from "@api/user.api";
import { useQuery } from "@tanstack/react-query";
import { IUser } from "types/user/user.type";

export const useGetUser = (id: string) => {
    return useQuery<IUser>({
        queryKey: ["user", id],
        queryFn: () => getUserById(id),
        enabled: !!id
    });
};

export const useGetUserChat = (ids: string[]) => {
    return useQuery<IUser>({
        queryKey: ["user", ids],
        queryFn: () => getUserById(ids[0]),
        enabled: ids.length == 1,
    });
};