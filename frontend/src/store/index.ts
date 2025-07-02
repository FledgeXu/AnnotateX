import { createStore, persist } from "easy-peasy";
import { authModel } from "@/store/models/authModel";
import type { StoreModel } from "@/store/types";
import { projectsModel } from "./models/projectsModel";

export const store = createStore<StoreModel>(
  persist(
    {
      auth: authModel,
      projects: projectsModel
    },
    {
      allow: ["auth"],
      storage: localStorage,
    },
  ),
);
