import { createSSRApp, createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from './store'
import './assets/css/main.css'
const isServer = typeof window === "undefined";

const app = isServer ? createSSRApp(App) : createApp(App);

app.use(router);
app.use(store);
app.mount("#app");
