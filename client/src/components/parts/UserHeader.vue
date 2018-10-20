<template>
  <header>
    <h1>
      <span class="nickname">{{xuser.nickname}}</span><span class="name">{{xuser.name}}</span>
      ({{xuser.reportCount}}レポート)
    </h1>
  </header>
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
