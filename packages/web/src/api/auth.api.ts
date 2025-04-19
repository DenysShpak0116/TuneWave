import { $api } from "./base.api";

export const login = (email: string, password: string) => {
    return $api.post("/auth/login", { email, password });
};

export const register = (data: {
    username: string;
    email: string;
    password: string;
}) => {
    return $api.post("/auth/register", data);
};

