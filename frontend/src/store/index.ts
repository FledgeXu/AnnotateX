import { createStore } from "easy-peasy";
import { authModel } from "@/store/models/authModel";
import type { StoreModel } from "@/store/types";

export const store = createStore<StoreModel>({
    auth: authModel,
});
