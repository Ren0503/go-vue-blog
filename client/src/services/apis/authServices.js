import { AxiosService } from './axiosService';
import { UsersService } from '../base/usersService';
import { NotificationService } from '../base/notificationService';

const login = (loginDto) => {
    NotificationService.setIsLoading(true);
    return new Promise((resolve, reject) => {
        AxiosService.post('auth/login', loginDto).then(res => {
            NotificationService.setIsLoading(false);
            if (res.data.success) {
                NotificationService.showSuccessAlert('Logged in successfully');
                const user = res.data.user;
                user.token = res.data.token;
                AxiosService.setUser(user);
                UsersService.saveUser(user);
            }
            resolve(res.data);
        }).catch(err => {
            debugger;
            NotificationService.setIsLoading(false);
            NotificationService.showToastError(err.message);
            reject(err);
        })
    });
};


const register = (userObject) => {
    NotificationService.setIsLoading(true);
    return new Promise((resolve, reject) => {
        AxiosService.post('users', userObject).then(({data}) => {
            NotificationService.setIsLoading(false);
            NotificationService.showSuccessAlert('Registered successfully');
            resolve(data);
        }).catch(err => {
            debugger;
            NotificationService.showToastError(err.message);
            NotificationService.setIsLoading(false);
            reject(err);
        });
    });
};


export const AuthService = {
    login,
    register,
};
