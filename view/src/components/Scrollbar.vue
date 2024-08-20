<template>
  <div></div>
</template>
<script>
export default {
  name: "Scrollbar",
  props:{
    topSize: {
      type: Number,
      default: 30
    },
    method: {
      type: Function,
      default: undefined
    },
  },

  data(){
    return {
      num: 1,
    }
  },
  created() {
      // this.listenScroll()
    window.addEventListener("scroll", this.listenScroll)
  },
  methods: {
    listenScroll(){
      if ((window.scrollY <= this.topSize) && this.num == 1 && this.method !== undefined) {
        this.method()
      }
      this.num = 0
      let that = this
      // 延迟 500 毫秒执行
      setTimeout(function () {
        that.num = 1
      }, 500)
    }
  }
}
</script>