import { Action, action } from 'easy-peasy';

export interface AuthModel {
    token: string | null;
    isAuthenticated: boolean;

    login: Action<AuthModel, string>;
    logout: Action<AuthModel>;
    loadFromStorage: Action<AuthModel>;
}

export const authModel: AuthModel = {
    token: null,
    isAuthenticated: false,

    login: action((state, token) => {
        state.token = token;
        state.isAuthenticated = true;
        localStorage.setItem('jwt', token);
    }),

    logout: action((state) => {
        state.token = null;
        state.isAuthenticated = false;
        localStorage.removeItem('jwt');
    }),

    loadFromStorage: action((state) => {
        const token = localStorage.getItem('jwt');
        if (token) {
            state.token = token;
            state.isAuthenticated = true;
        }
    }),
};


