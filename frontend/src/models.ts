export interface Response<T = unknown> {
  code: number;
  message: string;
  data: T;
}

export interface LoginToken {
  token: string;
}
