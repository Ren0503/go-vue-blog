import axios from 'axios';

const axiosInstance = axios.create({
    baseURL: process.env.BASE_URL,
    responseType: 'json',
    responseEncoding: 'utf8'
});

let cachedUser = {};

export const setUser = (user) => {
    cachedUser = user;
};

axiosInstance.interceptors.request.use((request) => {
    if (cachedUser.username)
        request.headers.authorization = "Bearer " + cachedUser.token;

    return request;
}, function (error) {
    debugger;
    return Promise.reject(error);
});

axiosInstance.interceptors.response.use((response) => {
    return response;
}, function (error) {
    if (error.response.status === 401 || error.response.status === 403) {
        UsersService.logout();
    }
});

function fetchPage(path, query = {page: 1, page_size: 12}) {
    return axiosInstance.get(`${path}?page=${query.page}&page_size=${query.page_size}`)
}

function get(path) {
    return axiosInstance.get(path)
}

function post(path, data, headers = null) {
    if (headers == null)
        return axiosInstance.post(path, data);
    else
        return axiosInstance.post(path, data, headers);
}

function put(path) {}

function _delete(path) {}

export const AxiosService = {
    axiosInstance, get, post, put, delete: _delete, setUser, fetchPage
};