<template>
  <div>
    <div slot="header" class="header">
      <el-page-header @back="goBack" :content="(friend ? friend.nick_name : '')"></el-page-header>
    </div>

    <div class="main"
         v-loading="loading"
         element-loading-text="加载中"
         element-loading-spinner="el-icon-loading"
    >
      <div v-for="(message, index) in messages" :key="index"  :class="(message.user_id == parseInt(friendId) ? 'friend-message' : 'my-message')">
        <el-avatar src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png"></el-avatar>

        <div class="item message">
          <span v-if="message.media == 1">{{message.content}}</span>
          <img v-if="message.media == 2" :src="message.content">
          </div>
      </div>
    </div>
    <el-footer>
      <el-form inline>
        <i class="el-icon-eleme" @click="showEnjo"></i>
        <el-form-item>
            <el-input type="textarea" size="large" autosize v-model="message" label="请输入信息" placeholder="请输入信息" clearable></el-input>
        </el-form-item>
        <el-form-item style="padding-top:4px;">
            <el-button type="primary" size="small" @click="sendMessage">发送</el-button>
            <el-button type="success" size="small" icon="el-icon-circle-plus-outline"></el-button>
        </el-form-item>
        <el-form-item>
           
        </el-form-item>
      </el-form>
    </el-footer>
    <el-drawer
      :visible.sync="drawer"
      :modal="false"
      direction="btt"
      :before-close="drawerClose">
      <el-tabs v-model="activeName" :stretch="true">
        <el-tab-pane label="😀  表情" name="first">
          <span class="emoji" v-for="(emoji, index) in emojis" :key="index" @click="selectEmoji(emoji)">{{emoji}}</span>
        </el-tab-pane>
        <el-tab-pane label="动态表情" name="second">
          <img v-for="(img, index) in emojis_img" :key="index" :src="img" alt="" class="emoji-img" @click="selectEmojiImg(img)">
        </el-tab-pane>
      </el-tabs>
    </el-drawer>
    <Scrollbar :top-size="30" :method="testCall"></Scrollbar>
  </div>
</template>

<script>
import {initSocket, socketClient, sendMessage} from "@/api/socket";
import {getFriend} from "@/api/friends";
import {getMessages} from "@/api/message";
import Scrollbar from "@/components/Scrollbar";
const EMOJIS = [
        "😀", "😁", "😂", "🤣", "😃", "😄", "😅", "😆", "😉", "😊", "😋", "😎", "😍", "😘", "😗", "😙", "😚", "☺", "🙂", "🤗", "🤔", "😐", "😑", "😶", "🙄", "😏", "😣", "😥", "😮", "🤐", "😯", "😪", "😫", "😴", "😌", "😛", "😜", "😝", "🤤", "😒", "😓", "😔", "😕", "🙃", "🤑", "😲", "☹", "🙁", "😖", "😞", "😟", "😤", "😢", "😭", "😦", "😧", "😨", "😩", "😬", "😰", "😱", "😳", "😵", "😡", "😠", "😷", "🤒", "🤕", "🤢", "🤧", "😇", "🤠", "🤡", "🤥", "🤓", "😈", "👿", "👹", "👺", "💀", "👻", "👽", "🤖", "💩", "😺", "😸", "😹", "😻", "😼", "😽", "🙀", "😿", "😾", "🏻", "🏼", "🏽", "🏾", "🏿", "🗣", "👤","👥","👫","👬"
,"👭","👂","👂🏻","👂🏼","👂🏽","👂🏾","👂🏿","👃","👃🏻","👃🏼","👃🏽","👃🏾","👃🏿","👣","👀","👁","👅","👄","💋","👓","🕶","👔","👕","👖","👗","👘","👙","👚","👛","👜","👝","🎒","👞","👟","👠","👡","👢","👑","👒","🎩","🎓","⛑","💄","💍","🌂","💼"
      ];
const EMOJIS_IMG = ["http://q.qqbiaoqing.com/q/2013/05/15/27229e5a8dddd37ceedec2f42ca2fa5f.gif", "http://q.qqbiaoqing.com/q/2010-5-4/ecf2ae2b0323d25c9d673e8218a2d85d.gif", "http://q.qqbiaoqing.com/q/2014/09/09/00625d3aebb803fe63636c783dc54183.gif", "http://q.qqbiaoqing.com/q/2014/07/22/e1c46311e324e53fddb32c90406b7141.gif", "http://q.qqbiaoqing.com/q/2014/09/11/d89047bef51bb6677f50887701f724f1.gif", "http://q.qqbiaoqing.com/q/2014/10/20/ba281059eb6de8f47bd10bc22523c3bf.gif", "http://q.qqbiaoqing.com/q/2013/01/26/0e548816baff4734fbe01579d84991ca.gif", "http://q.qqbiaoqing.com/q/2013/07/29/9a42d88a10b89e8b6274e337b8d61b7c.gif", "http://q.qqbiaoqing.com/q/2014/09/16/958102b66c78d70e505d4a13427f19cd.gif", "http://q.qqbiaoqing.com/q/2015/02/13/5f1de44d9d2a54d16ee4ffab4efb100c.gif", "http://q.qqbiaoqing.com/q/2013/08/31/b66973d0309d79474fede48b509cb629.gif", "http://q.qqbiaoqing.com/q/2014/01/29/0abd0e7cbe01788a4ccb645807ce4f22.gif", "http://q.qqbiaoqing.com/q/2015/09/21/199cb92be692dd64dcb7e5fbf21a73d4.gif", "http://q.qqbiaoqing.com/q/2013/10/17/462552706136fce42cb11ec8aa0018bd.gif", "http://q.qqbiaoqing.com/q/2017/06/19/ba841af482ac6423293e296568c7ffbd.gif", "http://q.qqbiaoqing.com/q/2015/06/19/ca0fdb73b7ed2aebf8cc3892f24f7828.gif", "http://q.qqbiaoqing.com/q/2010-8-6/7d121e162b589835ac48004c3c0dfc12.gif", "http://q.qqbiaoqing.com/q/2014/08/18/0ba0d9e24a6dd2338c91994d9859b640.gif", "http://q.qqbiaoqing.com/q/2014/12/01/e1516bef56d89fd540c8f720f2a3f14b.gif", "http://q.qqbiaoqing.com/q/2015/05/11/7b30dd8ee15826e0b2aba352e9a5a919.gif", "http://q.qqbiaoqing.com/q/2014/08/26/a012104ced4186cf6b3393191aa9015e.gif", "http://q.qqbiaoqing.com/q/2011-2-15/044e375e49c55e28d5100e16b05aab61.gif", "http://q.qqbiaoqing.com/q/2016/04/06/d7f480575e38a7066bd9dadd7540275d.gif", "http://q.qqbiaoqing.com/q/2014/08/19/f5e28aea3e8f83806014caa6275992bb.gif", "http://q.qqbiaoqing.com/q/2015/02/13/b1e8e1eb0402ca62cedec48c154be483.gif", "http://q.qqbiaoqing.com/q/2015/02/13/a5381cb4d73533b7e04b1ce2e15e9e11.gif", "http://q.qqbiaoqing.com/q/2017/01/09/91b566f72bbeb4ccf7bcdf89785573e2.gif", "http://q.qqbiaoqing.com/q/2014/03/28/9983abf09d1b00e3a622c22ef952b819.gif"]
export default {
  name: "Chat",
  components: {
    Scrollbar
  },
  data(){
    return {
      loading: true,
      friendId: 0,
      friend: null,
      message: "",
      socket: undefined,
      minId: 0,
      messages: [],
      drawer: false,
      activeName: 'first',
      emojis: Object.assign([], EMOJIS),
      emojis_img: Object.assign([], EMOJIS_IMG)
    }
  },
  created() {
    let friendId = this.$route.query.friendId
    this.friendId = friendId
    // initSocket
    this.getFriend(this.friendId)
    initSocket()
    this.socket = socketClient()
    this.socket.onmessage = this.recvMessage
    this.getMessages()
  },
  watch: {},
  methods: {
    testCall(){
      this.getMessages()
    },
    drawerClose(done){
      done()
    },
    goBack(){
      this.$router.push({path: "/home"})
    },
    recvMessage(event){
      let message = JSON.parse(event.data)
      console.log(message)
      this.messages = this.messages.concat(message)
    },
    getFriend(){
      getFriend(this.friendId).then(response => {
        if (response.code != 0) {
          this.$message({
            message: response.message,
            type: "error"
          })
          return
        }
        const { data } = response
        this.friend = data
      })
    },
    sendMessage(){
      if (this.message == "") {
        return;
      }
      let user = this.$store.getters.user
      let data = {
        content: this.message,
        to_id: parseInt(this.friendId),
        user_id: user ? user.id : '',
        cmd: 1, // 私聊
        media: 1, // 文本消息
      }
      this.messages = this.messages.concat(data)
      sendMessage(JSON.stringify(data))
      this.message = ""
    },
    getMessages(){
      getMessages(this.friendId, this.minId).then(response => {
        if (response.code != 0) {
          this.$message({
            message: response.message,
            type: "error"
          })
          return
        }
        let data = response.data
        let messagesTmp = [];
        for (let item of data) {
          this.minId = item.id
          messagesTmp.unshift(item)
        }
        this.messages = messagesTmp.concat(this.messages);
        // console.log(this.messages)
        this.loading = false
      })
    },
    showEnjo(){
      this.drawer = true
    },
    selectEmoji(emoji){
      this.message += emoji + ' '
    },
    selectEmojiImg(imgSrc){
      if (imgSrc == "") {
        return
      }
      let user = this.$store.getters.user
      let data = {
        content: imgSrc,
        to_id: parseInt(this.friendId),
        user_id: user ? user.id : '',
        cmd: 1, // 私聊
        media: 2, // 图片类型
      }
      this.messages.push(data)
      sendMessage(JSON.stringify(data))
      this.drawer = false
    }
  },
}
</script>

<style lang="scss" scoped>
.main {
  margin: 20px 10px;
  padding-bottom: 80px;
}
.header {
  position: fixed;
  top: 0;
  left: 0;
  z-index: 9000;
  height: 40px;
  line-height: 40px;
  background: #ffffff;
  width: 100%;
  padding: 10px 10px;
}
.el-footer {
  padding: 0;
  position: fixed;
  width: 100%;
  bottom: 0;
  left: 0;
  min-height: 40px;
  height: auto !important;
  border-top: 1px solid #B3C0D1;
  background-color: #ffffff;
  .el-form {
    margin-top: 10px;
  }
  .el-button {
    // margin-top: 7px;
  }
}

.friend-message {
  position: relative;
  margin-top: 10px;
  text-align: left;
  .item {
    height: auto;
    line-height: 30px;
    border: 1px solid lightgreen;
    padding: 5px 10px;
    word-break: break-all;
    text-align: left;
    font-size: 16px;
    font-weight: 600;
    background-color: lightgreen;
    margin-left: 45px;
    border-radius: 5px;
    display: inline-block;
  }
  .el-avatar {
    position: absolute;
    left: 0px;
  }
}

.my-message {
  position: relative;
  margin-top: 10px;
  text-align: right;
  .item {
    height: auto;
    line-height: 30px;
    border: 1px solid lightblue;
    padding: 5px 10px 0;
    word-break: break-all;
    text-align: left;
    font-size: 16px;
    font-weight: 600;
    background-color: lightblue;
    margin-right: 45px;
    border-radius: 5px;
    display: inline-block;
  }
  .el-avatar {
    position: absolute;
    right: 0px;
  }
}

.el-icon-eleme {
    height: 40px;
    font-size: 30px;
    margin-top: 8px;
    margin-right: 5px;
    color: #f6d84a;
}

.emoji {
  font-size: 30px; margin:15px 10px;
}
.emoji-img {
  width:50px;
  height: 50px;
  margin: 10px;
}
</style>
<style lang="scss">
.el-drawer { 
  height: 45% !important;
  &__header{
    align-items: center;
    color: #72767b;
    display: flex;
    margin-bottom: 0 !important;
    padding: 20px 20px 0;
  }
  &__body {
  padding: 5px 10px;
  text-align: left;
  flex: 1;
  overflow-y: scroll;
  }
}
</style>