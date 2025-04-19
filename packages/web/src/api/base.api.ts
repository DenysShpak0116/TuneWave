import axios from "axios";

export const $api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    headers: {
        "Content-Type": "application/json",
    },
});

let isRefreshing = false;
let failedQueue: any[] = [];

const processQueue = (error: any, token: string | null = null) => {
    failedQueue.forEach(prom => {
        if (error) {
            prom.reject(error);
        } else {
            prom.resolve(token);
        }
    });
    failedQueue = [];
};

$api.interceptors.response.use(
    (response) => response,
    async (error) => {
        const originalRequest = error.config;

        if (error.response?.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;

            if (!isRefreshing) {
                isRefreshing = true;
                try {
                    const res = await axios.post(`${import.meta.env.VITE_API_URL}/auth/refresh`, {}, { withCredentials: true });
                    const newAccessToken = res.data.accessToken;

                    $api.defaults.headers.common["Authorization"] = `Bearer ${newAccessToken}`;
                    processQueue(null, newAccessToken);
                    return $api(originalRequest);
                } catch (err) {
                    processQueue(err, null);
                    return Promise.reject(err);
                } finally {
                    isRefreshing = false;
                }
            }

            return new Promise(function (resolve, reject) {
                failedQueue.push({
                    resolve: (token: string) => {
                        originalRequest.headers["Authorization"] = `Bearer ${token}`;
                        resolve($api(originalRequest));
                    },
                    reject: (err: any) => reject(err),
                });
            });
        }

        return Promise.reject(error);
    }
);
