<style scoped>
  .alert {
    width: 100%;
    border-radius: 5px;
    border: 1px solid;
    border-color: #ff5656;
    background-color: #ffbcbc;
  }

  .alert p {
    display: inline-block;
    margin: 0;
    padding: 5px 10px;
  }

  .alert .close {
    float: right;
    padding: 5px 10px;
  }
</style>

<template lang="pug">
  div.alert(v-if="text")
    p {{text}}
    div.close(@click="this.close") X
</template>

<script>
  export default {
    props: {
      text: ""
    },
    methods: {
      close() {
        this.$emit("close")
      }
    },
    data() {
      return {}
    }
  }
</script>
