import request from "@/util/request"

/**
 * 添加好友
 * @param friendId
 */
export function addFriend(friendId){
    return request({
        method: 'post',
        url: '/friend',
        data: {
            friend_id: friendId
        }
    })
}

/**
 * 删除好友
 * @param friendId
 * @returns {AxiosPromise}
 */
export function deleteFriend(friendId){
    return request({
        method: 'delete',
        url: '/delete/friend?friendId=' + friendId,
    })
}
