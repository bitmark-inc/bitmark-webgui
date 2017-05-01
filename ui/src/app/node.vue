<style scoped>
  .action {
    float: right;
  }

  h3 .action .btn {
    border: none;
    margin-right: 10px;
    background: none;
    color: rgb(0, 96, 242);
    text-transform: uppercase;
    font-size: 16px;
    font-weight: bold;
    text-decoration: none;
  }

  h3 .action .btn:hover {
    color: rgb(126, 211, 33);
  }

  h3 .action .btn.stop:hover {
    color: red;
  }

  h3 .action .btn[disabled],
  h3 .action .btn[disabled]:hover {
    color: rgb(193, 193, 193);
    cursor: not-allowed;
  }
</style>

<template lang="pug">
  div
    h3 current chain
    p {{this.network}}
    h3 bitmark node
      div.action
        button.btn(@click="this.startBitmarkd", :disabled="this.bitmarkd.status==='started'") Start
        button.btn.stop(disabled, @click="this.stopBitmarkd", :disabled="this.bitmarkd.status==='stopped'") Stop
        router-link(tag="button", class="btn",to="/node/config") Config
    h3 prooferd node
      div.action
        button.btn(@click="this.startProoferd", :disabled="this.prooferd.status==='started'") Start
        button.btn.stop(disabled, @click="this.stopProoferd", :disabled="this.prooferd.status==='stopped'") Stop
        router-link(tag="button", class="btn", to="/node/config") Config
    h3 configuration
</template>

<script>
  import {
    getCookie
  } from "../utils"
  import axios from "axios"
  export default {
    methods: {
      startBitmarkd(e) {
        e.preventDefault();
        console.log("start")
        axios.post("/api/" + "bitmarkd", {
          option: "start"
        })
      },

      stopBitmarkd(e) {
        e.preventDefault();
        axios.post("/api/" + "bitmarkd", {
          option: "stop"
        })
      },

      startProoferd(e) {
        e.preventDefault();
        e.preventDefault();
        axios.post("/api/" + "prooferd", {
          option: "start"
        })
      },

      stopProoferd(e) {
        e.preventDefault();
        e.preventDefault();
        axios.post("/api/" + "prooferd", {
          option: "stop"
        })
      },

      fetchStatus(serviceName) {
        let service = this[serviceName]
        if (service.querying) {
          return
        }
        service.querying = true
        axios.post("/api/" + serviceName, {
            option: "status"
          })
          .then((resp) => {
            if (resp.data.ok) {
              service.status = resp.data.result
            } else {
              this.$emit("error", resp.data.result)
            }
          }).catch((e) => {
            this.$emit("error", e)
          })
          .then(() => {
            service.querying = false
          })
      }
    },

    mounted() {
      let network = getCookie("bitmark-webgui-network")
      if (!network) {
        this.$router.push("/chain")
      }
      this.network = network;
      this.bitmarkdTask = setInterval(() => {
        this.fetchStatus('bitmarkd')
      }, 2000)
      this.prooferdTask = setInterval(() => {
        this.fetchStatus('prooferd')
      }, 2000)
    },

    destroyed() {
      clearInterval(this.bitmarkdTask)
      clearInterval(this.prooferdTask)
    },

    data() {
      return {
        network: "",
        bitmarkdTask: null,
        prooferdTask: null,
        bitmarkd: {
          querying: false,
          status: "stopped"
        },
        prooferd: {
          querying: false,
          status: "stopped"
        }
      }
    }
  }
</script>
