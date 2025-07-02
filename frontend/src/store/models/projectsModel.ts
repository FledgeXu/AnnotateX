import type { Project } from "@/models";
import { action } from "easy-peasy";
import type { Action } from "easy-peasy";

export interface ProjectsModel {
    projects: Project[];
    addProject: Action<ProjectsModel, Project>;
    updateProjects: Action<ProjectsModel, Project[]>;
}


export const projectsModel: ProjectsModel = {
    projects: [],
    addProject: action((state, project) => {
        state.projects.push(project)
    }),
    updateProjects: action((state, projects) => {
        state.projects = projects
    }),
}
