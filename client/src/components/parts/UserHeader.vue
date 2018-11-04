<template>
  <div>
    <h1>
      <span class="nickname">{{xuser.nickname}}</span><span class="name">&lt;{{xuser.name}}&gt;</span>
      ({{xuser.reportCount}}レポート)
    </h1>
  </div>
</template>

<script>
export default {
  name: 'UserHeader',
  props: {
    xuser: {}
  }
}
</script>

<style lang="scss" scoped>
h1 {
  .nickname {
    font-size: 110%;
  }

  .name {
    color: #999;
  }
}

h2 {
  text-align: right;
  font-weight: normal;
}
</style>
