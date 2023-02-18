import {createRouter, createMemoryHistory, createWebHistory} from 'vue-router';
import Home from '../views/Home.vue';
import Config from '../views/Config.vue';

const isServer = typeof window === 'undefined';
const history = isServer ? createMemoryHistory() : createWebHistory();

const routes = [
    {
        path: '/',
        name: 'Home',
        component: Home
    },
    {
        path: '/config',
        name: 'Config',
        component: Config
    }
];

const router = createRouter({
    history,
    routes
});

export default router;
