import { AxiosService } from "../apis/axiosService";
import { LocalStorageService } from "./coreServices";


let observers = [];
const USER_KEY = 'user';


function hasRole(user, roleName) {
    if (!user)
        user = getUser();

    if (user == null)
        return false;

    return user.roles && user.roles.findIndex(role => role === roleName) !== -1;
}

function isAdmin(user) {
    return hasRole(user, 'ROLE_ADMIN');
}

function isNotAdmin() {
    return !isAdmin();
}

function isAuthor(user) {
    return hasRole(user, 'ROLE_AUTHOR')
}

function addUserExtraFields(user) {
    user.isAdmin = isAdmin(user);
    user.isAuthor = isAuthor(user);
    return user;
}

function notifyObservers(user) {
    // const cart = JSON.parse(StorageService.get(CART_KEY));
    observers.forEach(o => {
        o(user);
    });
}

const subscribe = (observer, receiveFirstState = true) => {
    // no more than one subscription
    if (!observers.includes(observer)) {
        observers.push(observer);
        if (receiveFirstState) {
            const user = JSON.parse(LocalStorageService.get(USER_KEY)) || {};
            addUserExtraFields(user);
            observer(user);
        }
    }
};

const unsubscribe = (observer) => {
    const index = observers.indexOf(observer);
    if (index > -1)
        observers.splice(index, 1);
};

const logout = () => {
    return new Promise((resolve, reject) => {
        debugger;
        LocalStorageService.clear(USER_KEY);
        notifyObservers({});
        resolve({success: true});
    });

};

const init = () => {
    const user = getUser();
    if (user && user.username) {
        AxiosService.setUser(user);
        addUserExtraFields(user);
    }
    return user;
};
const getUser = () => {
    return JSON.parse(LocalStorageService.get(USER_KEY));
};

const saveUser = (user) => {
    LocalStorageService.set(USER_KEY, JSON.stringify(user));
    addUserExtraFields(user);
    notifyObservers(user);
};
export const UsersService = {
    logout, init, getUser, subscribe, unsubscribe, saveUser
};
