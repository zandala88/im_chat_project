import request from '@/util/request'

/**
 * 获取用户锁甲的群组
 * @returns {AxiosPromise}
 */
export function getCommunities(){
    return request({
        url: '/communities',
        method: 'get'
    })
}

/**
 * 创建群组
 * @param data
 */
export function createCommunity(data){
    return request({
        url: '/community',
        method: 'post',
        data
    })
}

/**
 * 加入群聊
 * @param id
 */
export function addCommunity(id){
    return request({
        url : `/join/community/${id}`,
        method: 'post',
    })
}

/**
 * 获取群组信息
 * @param id
 */
export function getCommunityInfo(id){
    return request({
        url: '/community/'+id,
        method: 'get',
    })
}

