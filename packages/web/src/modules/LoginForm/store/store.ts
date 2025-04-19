import { create } from "zustand";
import { persist } from "zustand/middleware";
import { IUser } from "types/user/user.type";

interface AuthState {
    accessToken: string | null;
    user: IUser | null;
    setAccessToken: (token: string) => void;
    setUser: (user: IUser) => void;
    logout: () => void;
    isAuth: () => boolean;
}

export const useAuthStore = create<AuthState>()(
    persist(
        (set, get) => ({
            accessToken: null,
            user: null,
            setAccessToken: (token) => set({ accessToken: token }),
            setUser: (user) => set({ user }),
            logout: () => set({ accessToken: null, user: null }),
            isAuth: () => !!get().accessToken,
        }),
        {
            name: "auth-storage",
        }
    )
);
