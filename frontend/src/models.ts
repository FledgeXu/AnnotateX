export interface Response<T = unknown> {
  code: number;
  message: string;
  data: T;
}

export interface LoginToken {
  token: string;
}

export interface ProjectResponse {
  results: Project[]
}

export type ProjectStatus = "active" | "inactive" | "archived";


export interface Project {
  id: number;
  name: string;
  modality: string;
  status: ProjectStatus;
  description: string;
  createdAt: string;
  updatedAt: string;
}
