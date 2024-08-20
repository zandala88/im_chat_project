<template>
  <div class="main">
    <div slot="header" class="header">
      <el-page-header @back="goBack" content="我的群聊"></el-page-header>
    </div>
    <div class="group">
      <div class="list" v-for="(item, index) in groupList" :key="index" @click="groupChat(item.id)">
        <el-avatar src="http://q.qqbiaoqing.com/q/2013/05/15/27229e5a8dddd37ceedec2f42ca2fa5f.gif"></el-avatar>
        <div class="info">{{item.name}} </div>
        <el-badge :value="12" class="item" type="success"></el-badge>
      </div>
    </div>
  </div>
</template>
<script>
import {getCommunities} from "@/api/community"
export default {
  name: "CommunityList",
  data(){
    return {
      groupList:[]
    }
  },
  created() {
    this.geCommunities()
  },
  methods: {
    goBack(){
      this.$router.push({
        path: "/home"
      })
    },
    groupChat(groupId){
      this.$router.push({
        path: `/community/${groupId}`,
      })
    },
    geCommunities(){
      getCommunities().then(response => {
        console.log(response)
        if (response.code != 0) {
          this.$message({
            message: response.message,
            type: "error"
          })
          return
        }
        const { data }  = response
        this.groupList = data
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.main {
  margin: 20px 10px;
  padding-bottom: 80px;
}

.group {
  margin-top: 20px;
  height: auto;
  .list {
    width: 100%;
    height: 40px;
    display: flex;
    border-bottom: 1px solid #dddddd;
    padding-bottom: 5px;
    margin-top: 10px;
    .info {
      line-height: 40px;
      height: 40px;
      font-size: 16px;
      font-weight: 600;
      margin-left: 10px;
    }
  }
}
</style>

<style lang="scss">
.el-page-header__content {
  font-size: 14px;
}
</style>