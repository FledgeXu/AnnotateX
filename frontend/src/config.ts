import axios from "axios";
import type { Store } from "easy-peasy";
import type { StoreModel } from "@/store/types";

export const createAPI = (store: Store<StoreModel>) => {
    const api = axios.create({
        baseURL: import.meta.env.VITE_BASE_URL,
        timeout: 360,
    });

    api.interceptors.request.use((config) => {
        const token = store.getState().auth.token;
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    })

    return api;
}
