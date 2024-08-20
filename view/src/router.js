import VueRouter from "vue-router";
import Login from "@/components/Login";
import Register from "@/components/Register";
import Home from "@/components/Home";
import store from "store";
import Friend from "@/components/Friends";
import My from "@/components/My";
import Chat from "@/components/Chat";
import vueStore from './store'
import Communities from "@/components/Communities";
import GroupChat from "@/components/GroupChat";
const routes = [
    {
        path: '/',
        redirect: ''
    },
    {
        path: '/home',
        component: Home
    },
    {
        path: "/login",
        name: "Login",
        component: Login
    },
    {
        path: '/register',
        name: 'Register',
        component: Register
    },
    {
        path: '/friend',
        name: "Friend",
        component: Friend
    },
    {
        path: '/communities',
        name: "Communities",
        component: Communities
    },
    {
        path: '/community/:groupId',
        name: "GroupChat",
        component: GroupChat
    },
    {
        path: '/my',
        name: "My",
        component: My
    },
    {
        path: "/chat",
        name: "Chat",
        component: Chat
    }
]

const router = new VueRouter({
    routes
})

// 路由白名单
let whiteRoutesPath = [
    "/login",
    "/register"
];
router.beforeEach(async(to, from, next) => {
    // 获取用户的基本信息以及 token
    let user = store.get("token")
    let token = store.get("token") || (user && user.token)
    console.log(to, from)
    if ((token == "" || token == undefined) && (whiteRoutesPath.indexOf(to.path) === -1)) {
        if (to.path == "/register") {
            next()
        } else {
            next({ path: "/login" })
        }
    } else if ((token !== "" && token !== undefined) && (whiteRoutesPath.indexOf(to.path) >= 0)) {
        next({ path: "/home" })
    } else if (to.path == "/") { 
        next({path: "/home"})
    }
    // 是否存在 state
    let storeUser = vueStore.getters.user || user
    vueStore.commit("setUser", storeUser)
    next()
})


export default router