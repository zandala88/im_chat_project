import request from "@/util/request"

/**
 * 获取用户的聊天记录
 * @param friendId
 * @param maxId
 */
export function getMessages(friendId, maxId){
    return request({
        url : `/messages/${friendId}`,
        method: "get",
        params: {
            maxId: maxId
        }
    })
}

// 获取群组消息
export function getCommunityMessages(communityId){
    return request({
        url : `/community_messages/${communityId}`,
        method: "get",
    })
}