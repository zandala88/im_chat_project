import request from "@/util/request"

/**
 * 获取当前用户的好友
 */
export function getFriends(){
    return request({
        url: '/friends',
        method: 'get',
    })
}

/**
 * 获取好友信息
 * @param friendId
 */
export function getFriend(friendId){
    return request({
        url : '/friend/'+friendId,
        method: 'get'
    })
}