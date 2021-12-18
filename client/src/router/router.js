import Vue from 'vue';
import Router from 'vue-router';
import ArticleList from '../components/articles/ArticleList';
import ArticleDetails from '../components/articles/ArticleDetails';
import Login from '../components/auth/Login';
import Register from '../components/auth/Register';
import Logout from '../components/auth/Logout';
import ArticleListByTag from "@/components/articles/ArticleListByTag";
import ArticleListByCategory from "@/components/articles/ArticleListByCategory";
import ArticleEdit from "@/components/articles/ArticleEdit";

Vue.use(Router);

export const router = new Router({
    mode: 'history',
    hash: false,

    routes: [
        {
            path: '',
            redirect: '/articles'
        },
        {
            path: '/articles/new',
            exact: true,
            name: 'article-new',
            props: true,
            component: ArticleEdit,
        },
        {
            path: '/articles/:slug',
            exact: true,
            name: 'article-details',
            props: true,
            component: ArticleDetails,
        },
        {
            path: '/articles',
            name: 'article-list',
            exact: true,
            component: ArticleList
        },
        {
            path: '/articles/by_tag/:slug',
            exact: true,
            name: 'article-list-by-tag',
            component: ArticleListByTag,
        },
        {
            path: '/articles/by_category/:slug',
            exact: true,
            name: 'article-list-by-category',
            component: ArticleListByCategory,
        },
        {
            path: '/login', 
            component: Login, 
            name: 'login', 
            onlyGuest: true
        },
        {
            path: '/register', 
            component: Register, 
            name: 'register', 
            onlyGuest: true
        },
        {
            path: '/logout', 
            component: Logout, 
            name: 'logout'
        },
    ],
});
