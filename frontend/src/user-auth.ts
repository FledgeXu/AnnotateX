import { HttpStatusCode } from 'axios';
import { createAPI } from "@/config"
import { store } from "@/store"
import { redirect } from "@tanstack/react-router";


export const useUserAuth = () => {
    const api = createAPI(store);;

    const isLogin = async () => {
        try {
            const res = await api.get("/v1/users/me")
            return res.status = HttpStatusCode.Ok;
        } catch (err) {
            return false;
        }
    };

    const check = async () => {
        if (!(await isLogin())) {
            throw redirect({ to: "/login" });
        }
    }
    return { check }
}

export type IUserAuth = ReturnType<typeof useUserAuth>;
