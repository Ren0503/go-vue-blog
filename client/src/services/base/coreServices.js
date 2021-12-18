export const LocalStorageService = {
    get(key) {
        return localStorage.getItem(key);
    },
    set(key, value) {
        localStorage.setItem(key, value);
    },
    remove(key) {
        localStorage.removeItem(key);
    },
    clear(key) {
        localStorage.removeItem(key);
    },
    delete(key) {
        localStorage.removeItem(key);
    },
    removeItem(key) {
        localStorage.removeItem(key);
    },
    deleteItem(key) {
        localStorage.removeItem(key);
    },
};

export const SessionStorageService = {
    get(key) {
        return sessionStorage.getItem(key);
    },
    set(key, value) {
        sessionStorage.setItem(key, value);
    },
    remove(key) {
        sessionStorage.removeItem(key);
    },
    clear(key) {
        sessionStorage.removeItem(key);
    },
    delete(key) {
        sessionStorage.removeItem(key);
    },
    clearAll() {
        sessionStorage.clear();
    }
};