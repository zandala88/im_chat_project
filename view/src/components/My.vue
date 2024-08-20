<template>
  <div>
    <div slot="header" class="clearfix">
      <el-page-header @back="goBack" content="个人中心"></el-page-header>
    </div>
    <el-main>
      <el-row>
        <el-col :span="24" class="user-info">
          <el-avatar src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"></el-avatar>
          <div class="mate">{{user.nick_name}} <span>ID: {{user.id}}</span></div>
        </el-col>
        <el-col :span="24">
          通知公告
        </el-col>
        <el-col :span="24">
          <el-button type="text" @click="addFriend">添加好友</el-button>
        </el-col>
        <el-col :span="24">
          <el-button type="text" @click="addCommunity">加入社群</el-button>
        </el-col>
        <el-col :span="24">
          <el-button type="text" @click="showCommunity">创建社群</el-button>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span="24">
          获取代码
        </el-col>
      </el-row>
    </el-main>
    <el-button @click="logout">退出登录</el-button>
  </div>
</template>

<script>
import {addFriend} from "@/api/user"
import {doLogout} from "@/api/login"
import {createCommunity, addCommunity} from "@/api/community";
export default {
  name: "My",
  data(){
    return {
      user: undefined
    }
  },
  created() {
    this.user = this.$store.getters.user
  },
  methods: {
    goBack(){
      this.$router.push({path:'/home'})
    },
    addFriend: function () {
      this.$prompt("请输入用户 ID", "提示", {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /\d+/,
        inputErrorMessage: 'ID格式不正确',
        center: true
      }).then(({value}) => {
        console.log(value)
        addFriend(value).then(response => {
          if (response.code != 0) {
            this.$message({
              message: response.message,
              type: "error",
              center:true
            })
            return;
          }
          this.$message({
            message: "添加好友成功",
            type: "success"
          })
          this.$router.push({path: "/friend"})
        })
      })
    },
    logout(){
      doLogout()
      this.$router.push({path: "/login"})
    },
    showCommunity(){
      let user = this.$store.getters.user
      let that = this
      // 创建社群
      this.$prompt("请输入群名称", "提示", {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /\S{3,25}/,
        inputErrorMessage: '群名称格式不正确',
        center: true
      }).then((data) => {
        that.createCommunity({
          name: data.value,
          ownerId:user.id
        })
      })
    },
    createCommunity(data){
      createCommunity(data).then((response) => {
        if (response.code != 0) {
          this.$message({
            message: response.message,
            type: 'error'
          })
          return
        }

        this.$message({
          message: '创建成功',
          type: 'primary'
        })
        this.router.push({
          path: '/communities'
        })
      })
    },
    addCommunity(){
      let that = this
      this.$prompt("请输入群ID", "提示", {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /\d+/,
        inputErrorMessage: '群ID格式不正确',
        center: true
      }).then((data) => {
        addCommunity(data.value).then(response => {
          if (response.code != 0) {
            that.$message({
              message: response.message,
              type: "error"
            })
            return
          }
          that.$message({
            message: "加群成功",
            type: "success"
          })
          that.$router.push({
            path: `/community/${data.value}`
          })
        })
      })
    }
  }
}
</script>

<style lang="scss" scoped>
.el-header {
  background-color: #B3C0D1;
  color: #333;
  text-align: center;
  line-height: 60px;
}
.el-col {
  border-bottom: 1px solid #B3C0D1;
  margin: 5px;
  /*padding: 5px;*/
  line-height: 50px;
  height: 50px;
  text-align: left;
}
.el-message-box {
  width: 100%;
}

.user-info {
  height: 50px;
  line-height: 50px;
  display: flex;
  margin-bottom: 5px;
  .mate {
    font-size: 16px;
    font-weight: 600;
    display: inline-block;
    span {
      color: #dddddd;
      font-size: 14px;
      margin-left: 10px;
    }
  }
}
</style>

<style lang="css" scoped>
.el-message-box {
  display: inline-block;
  padding-bottom: 10px;
  vertical-align: middle;
  background-color: #FFF;
  border-radius: 4px;
  border: 1px solid #EBEEF5;
  font-size: 18px;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,.1);
  text-align: left;
  overflow: hidden;
  -webkit-backface-visibility: hidden;
  backface-visibility: hidden;
  width: 100% !important;
}

</style>