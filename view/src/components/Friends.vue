<template>
  <div>
    <el-card class="box-card" body-style='{padding: 0}' shadow="never">
      <div slot="header" class="clearfix">
        <el-page-header @back="goBack" content="好友列表"></el-page-header>
      </div>
      <div v-for="friend in friends" :key="friend.id" class="text item" align="left">
        <div class="user-item" @click="goChat(friend.id)">
          <el-avatar src=""></el-avatar>
          <span class="text">{{friend.nick_name }}</span>
          <i class="el-icon-arrow-right"></i>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script>
import {getFriends} from "@/api/friends";

export default {
name: "Friend",
  created() {
    this.getFriends()
  },
  data() {
    return {
      friends: [],
    }
  },
  methods: {
    goChat(friendId){
      this.$router.push({path: "/chat", query:{friendId: friendId}})
    },
    goBack(){
      this.$router.push({path: "/home"})
    },
    getFriends(){
      const loading = this.$loading({
        lock: true,
        text: '好友加载中...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      });
      getFriends().then(response => {
        loading.close()
        if (response.code != 0) {
          this.$message({
            message: response.message,
            type: "error"
          })
          return
        }
        this.friends = response.data
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.user-item {
  display: block;
  height: 60px;
  line-height: 60px;
  border-bottom: 1px solid #B3C0D1;
  position: relative;
  vertical-align:middle;
  .text {
    margin-left: 50px;
  }
  .el-avatar {
    position: absolute;
    top: 10px;
  }
  i {
    position: absolute;
    right: 5px;
    top: 20px;
  }
}
</style>