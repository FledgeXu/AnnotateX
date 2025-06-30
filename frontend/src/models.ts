export interface Response<T = unknown> {
  code: number;
  message: string;
  data: T;
}

export interface LoginToken {
  token: string;
}

export type ProjectStatus = "active" | "inactive" | "archived";

export interface Project {
  id: number;
  name: string;
  modality: string;
  status: ProjectStatus;
  description: string;
  created_at: string;
  updated_at: string;
}
