import { HttpStatusCode } from 'axios';
import { createAPI } from "@/config"
import { store } from "@/store"


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
    return { isLogin }
}

export type IUserAuth = ReturnType<typeof useUserAuth>;
