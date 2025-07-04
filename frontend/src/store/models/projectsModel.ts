import type { Project, ProjectSortMode } from "@/models";
import { action, computed } from "easy-peasy";
import Fuse from "fuse.js";
import type { Action, Computed } from "easy-peasy";


export interface ProjectsModel {
    projects: Project[];
    sortMode: ProjectSortMode;
    searchQuery: string;

    // computed
    sortedProjects: Computed<ProjectsModel, Project[]>;
    visibleProjects: Computed<ProjectsModel, Project[]>;

    // actions
    setSortMode: Action<ProjectsModel, ProjectSortMode>;
    setSearchQuery: Action<ProjectsModel, string>;
    updateProjects: Action<ProjectsModel, Project[]>;
    addProject: Action<ProjectsModel, Project>;
}

export const projectsModel: ProjectsModel = {
    projects: [],
    sortMode: "createTimeDesc",
    searchQuery: "",

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
                result.sort((a, b) => new Date(a.created_at).getTime() - new Date(b.created_at).getTime());
                break;
            case "createTimeDesc":
                result.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());
                break;
        }
        return result;
    }),

    visibleProjects: computed((state) => {
        const projects = state.sortedProjects;
        const query = state.searchQuery.trim();

        if (!query) return projects;

        const fuse = new Fuse(projects, {
            keys: ["name", "description"],
            threshold: 0.3,
        });

        return fuse.search(query).map(result => result.item);
    }),

    setSortMode: action((state, mode) => {
        state.sortMode = mode;
    }),

    setSearchQuery: action((state, query) => {
        state.searchQuery = query;
    }),

    updateProjects: action((state, projects) => {
        state.projects = projects;
    }),

    addProject: action((state, project) => {
        state.projects.push(project);
    }),
};

