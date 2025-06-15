import { createStore } from 'easy-peasy';
import { authModel } from './models/authModel';
import { StoreModel } from './types';

export const store = createStore<StoreModel>({
    auth: authModel,
});
