<template>
  <div class="overlay" v-if="visible"></div>
</template>

<script>
export default {
  name: 'Overlay',
  props: {
    visible: Boolean
  }
}
</script>

<style scoped>
.overlay {
  z-inde: 9999;
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;

  background-color: rgba(0, 0, 0, 0.5);
}
</style>
