import axios from "axios"
import qs from "querystring"
import store from "store"
import Message from "element-ui/packages/message/src/main";
// import MessageBox from "element-ui/packages/message-box/src/main";
const REQUEST_BASE = "http://localhost:8080"

const service = axios.create({
    baseURL: REQUEST_BASE,
    timeout: 5000,
    // headers: {
    //     'Content-Type':'application/x-www-form-urlencoded'
    // }
})
// querystring
service.interceptors.request.use((config) => {
    let user = store.get("user")
    let token = store.get("token")
    let data = {
        user_id: user ? user.id : 0,
        token: token ? user.token : ''
    }
    //
    // if (config.data != undefined) {
    //     config.data = Object.assign(data, config.data)
    // }
    config.url += "?"+qs.stringify(data)
    return config
}, (error) => {
    return Promise.reject(error)
})

// response interceptor
service.interceptors.response.use(
    response => {
        const res = response.data
        return res
        // if the custom code is not 20000, it is judged as an error.
        // if (res.code !== 0) {
        //     Message({
        //         message: res.message || 'Error',
        //         type: 'error',
        //         duration: 5 * 1000
        //     })
        //
        //     // 50008: Illegal token; 50012: Other clients logged in; 50014: Token expired;
        //     if (res.code === 50008 || res.code === 50012 || res.code === 50014) {
        //         // to re-login
        //         MessageBox.confirm('You have been logged out, you can cancel to stay on this page, or log in again', 'Confirm logout', {
        //             confirmButtonText: 'Re-Login',
        //             cancelButtonText: 'Cancel',
        //             type: 'warning'
        //         }).then(() => {
        //
        //         })
        //     }
        //     return Promise.reject(new Error(res.message || 'Error'))
        // } else {
        //     return res
        // }
    },
    error => {
        console.log('err' + error) // for debug
        Message({
            message: error.message,
            type: 'error',
            duration: 5 * 1000
        })
        return Promise.reject(error)
    }
)

export default service