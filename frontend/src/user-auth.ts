import { HttpStatusCode } from "axios";
import { createAPI } from "@/config";
import { store } from "@/store";
import { toast } from "sonner";

export const useUserAuth = () => {
  const api = createAPI(store);

  const isLogin = async () => {
    try {
      const res = await api.get("/v1/users/me");
      return (res.status = HttpStatusCode.Ok);
    } catch {
      toast.error(`Failed to get user info.`);
      return false;
    }
  };
  return { isLogin };
};

export type IUserAuth = ReturnType<typeof useUserAuth>;
