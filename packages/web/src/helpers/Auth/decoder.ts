import { LoginResponse } from "@modules/LoginForm/types/loginResponse";

const decodeBase64Cookie = <T = string>(cookieValue: string): T | null => {
    try {
        const decoded = atob(cookieValue);
        return JSON.parse(decoded) as T;
    } catch (e) {
        console.error('Ошибка при декодировании:', e);
        return null;
    }
};

export const getUserFromCookie = (): LoginResponse | null => {
    const cookies = document.cookie.split('; ');
    const authCookie = cookies.find(row => row.startsWith('auth_data='));

    if (!authCookie) return null;

    const cookieValue = authCookie.split('=')[1];
    const userData = decodeBase64Cookie<LoginResponse>(cookieValue);
    document.cookie = 'auth_data=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT';

    if (userData?.user) {
        console.log("HEREE");
        return userData;
    }

    return null;
};