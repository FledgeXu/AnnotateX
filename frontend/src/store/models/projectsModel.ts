import type { Project, ProjectSortMode } from "@/models";
import { action, computed } from "easy-peasy";
import type { Action, Computed } from "easy-peasy";


export interface ProjectsModel {
    projects: Project[];
    sortMode: ProjectSortMode;
    // Computed
    orderBy: Computed<ProjectsModel, string>;
    order: Computed<ProjectsModel, string>;
    // actions
    setSortMode: Action<ProjectsModel, ProjectSortMode>;
    updateProjects: Action<ProjectsModel, Project[]>;
    clearProjects: Action<ProjectsModel, void>;
}

export const projectsModel: ProjectsModel = {
    projects: [],
    sortMode: "createTimeDesc",

    orderBy: computed((state) => {
        switch (state.sortMode) {
            case "createTimeDesc":
            case "createTimeAsc":
                return "created_at"
            case "nameAsc":
            case "nameDesc":
                return "name"
        }
    }),


    order: computed((state) => {
        switch (state.sortMode) {
            case "createTimeAsc":
            case "nameAsc":
                return "asc"
            case "createTimeDesc":
            case "nameDesc":
                return "desc"
        }
    }),

    setSortMode: action((state, mode) => {
        state.sortMode = mode;
    }),

    updateProjects: action((state, projects) => {
        state.projects = projects;
    }),

    clearProjects: action((state) => {
        state.projects = [];
    }),
};

