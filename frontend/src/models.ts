export interface Response<T = unknown> {
  code: number;
  message: string;
  data: T;
}

export interface Paged<T = unknown> {
  limit: number;
  offset: number;
  total: number;
  results: T[]
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

export const PROJECT_SORT_MODES = [
  "createTimeDesc",
  "createTimeAsc",
  "nameAsc",
  "nameDesc",
] as const;

export type ProjectSortMode = typeof PROJECT_SORT_MODES[number];

