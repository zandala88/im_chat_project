import request from "@/util/request"
import store from "store"
export function doLogin(data){
    return request({
        method: "POST",
        url: '/user/login',
        data
    })
}

export function doRegister(data){
    return request({
        url: '/user/register',
        method: 'post',
        data
    })
}

export function doLogout() {
    return store.clearAll()
    // return request({
    //     method: "GET",
    //     url: '/user/logout',
    //     params
    // })
}