import type { Project, ProjectSortMode } from "@/models";
import { action, computed } from "easy-peasy";
import type { Action, Computed } from "easy-peasy";


export interface ProjectsModel {
    projects: Project[];
    sortMode: ProjectSortMode;

    // computed
    sortedProjects: Computed<ProjectsModel, Project[]>;

    // actions
    addProject: Action<ProjectsModel, Project>;
    updateProjects: Action<ProjectsModel, Project[]>;
    setSortMode: Action<ProjectsModel, ProjectSortMode>;
}

export const projectsModel: ProjectsModel = {
    projects: [],
    sortMode: "createTimeDesc",

    sortedProjects: computed((state) => {
        const result = [...state.projects];

        switch (state.sortMode) {
            case "nameAsc":
                result.sort((a, b) => a.name.localeCompare(b.name));
                break;
            case "nameDesc":
                result.sort((a, b) => b.name.localeCompare(a.name));
                break;
            case "createTimeAsc":
                result.sort(
                    (a, b) =>
                        new Date(a.created_at).getTime() -
                        new Date(b.created_at).getTime()
                );
                break;
            case "createTimeDesc":
                result.sort(
                    (a, b) =>
                        new Date(b.created_at).getTime() -
                        new Date(a.created_at).getTime()
                );
                break;
        }

        return result;
    }),

    addProject: action((state, project) => {
        state.projects.push(project);
    }),

    updateProjects: action((state, projects) => {
        state.projects = projects;
    }),

    setSortMode: action((state, mode) => {
        state.sortMode = mode;
    }),
};
