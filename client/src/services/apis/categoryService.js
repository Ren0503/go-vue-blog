import { AxiosService } from './axiosService';
import { createPromiseWrapper } from '../../utils/promise';

export const CategoryService = {
    fetchAll() {
        return createPromiseWrapper(AxiosService.get('categories'), 'Categories fetched!');
    },
};
