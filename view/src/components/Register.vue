<template>
  <div>
    <h3 style="text-align: center">用户注册</h3>
    <el-form
        ref="logForm"
        :model="form"
    >
      <el-form-item label="用户名称">
        <el-input v-model="form.nickname" clearable placeholder="用户名称"></el-input>
      </el-form-item>
      <el-form-item label="手机号">
        <el-input v-model="form.mobile" clearable placeholder="手机号"></el-input>
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="form.passwd" show-password clearable placeholder="请输入密码"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="doRegister">注 册</el-button>
      </el-form-item>
    </el-form>
    <div style="text-align: center">
      <el-link type="primary" @click="redirect('login')">用户登录</el-link>
      |
      <el-link type="warning" @click="redirect('findPwd')">找回密码</el-link>
    </div>
  </div>
</template>

<script>
import {doRegister} from "@/api/login"
export default {
  name: 'Register',
  data: function () {
    return {
      form: {
        nickname: "",
        mobile: "",
        passwd: ""
      }
    }
  },
  methods: {
    redirect(path){
      this.$router.push(path)
    },
    doRegister(){
      doRegister(this.form).then((response) => {
        if (response.code != 0) {
          this.$message({
            type: 'error',
            message: response.message
          })
        } else {
          this.$message({
            type: 'primary',
            message: '注册成功'
          })
        }
      }).catch(err => {
        console.log(err)
      })
    }
  }
}
</script>
