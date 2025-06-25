export interface Response<T = any> {
    code: number;
    message: string;
    data: T;
}

export interface LoginToken {
    token: string,
}
