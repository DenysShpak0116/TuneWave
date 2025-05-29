import { $api } from "./base.api";

export const login = async (email: string, password: string) => {
    return await $api.post("/auth/login", { email, password });
};

export const register = async (data: {
    username: string;
    email: string;
    password: string;
}) => {
    return await $api.post("/auth/register", data);
};

export const sendTokenToEmail = async (email: string) => {
    return await $api.post("/auth/forgot-password", { email })
}

export const resetPassword = async (newPassword: string, token: string) => {
    return await $api.post("/auth/reset-password", { newPassword, token })
}

export const logout = async () => {
    return await $api.post("/auth/logout");
}