import store from "store"
import qs from "querystring"
let socket = null

// 初始化 socket
export function initSocket(){
    const user = store.get("user")
    const token = store.get("token")
    const data = {
        user_id: user.id,
        token: token
    }
    let queryStr = qs.stringify(data)
    const socketUrl = "ws://localhost:8080/chat?"+queryStr
    socket = new WebSocket(socketUrl)
    console.log(socket);

    socket.onopen = open
    socket.onerror = error
    // socket.onmessage = message
    socket.onclose = close
}

// open
function open(event){
    console.log(event);
}

function close(event){
    console.log(event);
    socket.close()
}

// function message(event){
//     console.log("收到了数据", event.data);
// }

// 发送 data
export function sendMessage(data){
    console.log(data)
    socket.send(data)
}

function error(event){
    console.log(event);
}

/**
 * 将 socket 注入到其他的页面里面
 * @returns {null}
 */
export function socketClient(){
    return socket
}
