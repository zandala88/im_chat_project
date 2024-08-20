<template>
  <el-container>
    <el-header>Golang IM</el-header>
    <el-main>
      <el-row>
        <el-col :span=24>
          <el-button icon="el-icon-chat-dot-round" type="primary" @click="redirect('friend')">
          好友
          </el-button>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span=24>
          <el-button icon="el-icon-chat-round" type="success" @click="redirect('communities')">
            群聊
          </el-button>
        </el-col>
      </el-row>
      <el-row>
        <el-col :span=24>
      <el-button icon="el-icon-s-custom" type="info"  @click="redirect('my')">
        我的
      </el-button>
        </el-col>
      </el-row>
    </el-main>
    <el-footer>
      <div>Copyright © 2020 </div>
      <div><a href="http://huike.fm126.top">http://huike.fm126.top</a></div>
    </el-footer>
  </el-container>
</template>

<script>
export default {
name: "Home",
  methods: {
    redirect(type) {
      this.$router.push({path: '/' + type})
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
.el-row {
  margin-top: 20px;
  .el-button {
    width: 100%;
  }
}

.el-footer {
  background-color: #B3C0D1;
  position: fixed;
  bottom: 0px;
  text-align: center;
  padding: 0;
  margin: 0;
  width: 100%;
  div {
    padding: 5px;
  }
}
</style>