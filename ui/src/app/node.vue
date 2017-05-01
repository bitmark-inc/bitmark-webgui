<style scoped>

</style>

<template lang="pug">
  div
    h3 current chain
    h3 bitmark node
    h3 prooferd node
    h3 configuration
</template>

<script>
import {getCookie} from "../utils"
  export default {
    mounted () {
      let network = getCookie("bitmark-webgui-network")
      if (!network) {
        this.$router.push("/chain")
      }
    },
    data() {
      return {}
    }
  }
</script>
