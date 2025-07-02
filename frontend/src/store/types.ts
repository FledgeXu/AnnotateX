import type { AuthModel } from "@/store/models/authModel";
import type { ProjectsModel } from "@/store/models/projectsModel";

export interface StoreModel {
  auth: AuthModel;
  projects: ProjectsModel
}
