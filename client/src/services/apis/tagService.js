import { AxiosService } from './axiosService';
import { createPromiseWrapper } from '../../utils/promise';

export const TagService = {
    fetchAll() {
        return createPromiseWrapper(AxiosService.get('tags'), 'Tags fetched!');
    },
}