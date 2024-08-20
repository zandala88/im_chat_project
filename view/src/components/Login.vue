<template>
  <div>
    <h3 style="text-align: center">用户登录</h3>
    <el-form
        ref="logForm"
        :model="form"
    >
      <el-form-item label="用户手机">
        <el-input v-model="form.mobile" clearable placeholder="用户手机"></el-input>
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="form.passwd" show-password clearable placeholder="请输入密码"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="doLogin">登 录</el-button>
      </el-form-item>
    </el-form>
    <div style="text-align: center">
      <el-link type="primary" @click="redirect('register')">用户注册</el-link>
      |
      <el-link type="warning" @click="redirect('findPwd')">找回密码</el-link>
    </div>
  </div>
</template>

<script>
import {doLogin} from "@/api/login"
import store from "store";
export default {
  name: 'Login',
  data: function () {
    return {
      form: {
        mobile: "",
        passwd: ""
      }
    }
  },
  methods: {
    redirect(path){
      this.$router.push(path)
    },
    doLogin(){
      doLogin(this.form).then((response) => {
        if (response.code != 0) {
          this.$message({
            type: 'error',
            message: response.message
          })
        } else {
          const { data } = response
          store.set('token', data.token)
          store.set("user", data)
          this.$store.commit("setUser", data)
          this.$message({
            type: 'primary',
            message: '登录成功'
          })
          this.$router.push({path: "/home"}, () => {})
        }
      }).catch(err => {
        console.log(err)
      })
    }
  }
}
</script>
