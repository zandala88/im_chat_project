import Vue from "vue";
import Vuex from 'vuex'
Vue.use(Vuex)
import store from "store"
const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
    state : {
        user: store.get("user"),
        token: store.get("token"),
    },
    getters: {
        user: state => {
            return state.user
        },
        token: state => {
            return state.token
        }
    },
    actions: {
        // userInfo({commit, state}){
        //     console.log(1111)
        //     console.log(commit, state)
        // }
    },
    mutations: {
        setUser(state, data){
            state.user = data || undefined
            state.token = data ? data.token : ''
        },
    },
    strict: debug,
})
