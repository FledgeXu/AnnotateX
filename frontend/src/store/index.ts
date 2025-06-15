import { createStore } from "easy-peasy";
import { authModel } from "./models/authModel";
import type { StoreModel } from "./types";

export const store = createStore<StoreModel>({
    auth: authModel,
});
