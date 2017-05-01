<style lang="" scoped>
  .panel {
    max-width: 730px;
    margin: 0 auto;
  }

  .box {
    border: none;
    background-color: rgb(237, 240, 244);
  }

  .panel-heading .sub {
    font-style: italic;
  }

  .option {
    font-weight: normal;
  }

  .option input {
    margin-right: 5px;
  }

  .option .help-text {
    font-size: 12px;
    margin: 4px 17px;
  }

  .panel-footer {
    border: none;
    background: none;
    text-align: right;
  }

  .start-node {
    display: inline-block;
    border: none;
    padding: 9px 24px;
    background-color: white;
    color: rgb(0, 96, 242);
  }

  .start-node:hover {
    color: black;
  }
</style>

<template lang="pug">
div
  h4 start bitmark node
  div.panel.panel-default.box
    div.row
      div.col-md-3
        div.panel-heading
          h5.title select chain
          p.sub Bitmark provide two diffrent chains to let the bitmarkd join in. They are testing, bitmark.
      div.col-md-9
        div.panel-body
          label.option
            input(type="radio")
            .
              TESTING
            p.help-text Link to public test bitmark network, to pay the transactions, please contact us.
          label.option
            input(type="radio")
            .
              BITMARK
            p.help-text Link to public bitmark network, pay the transactions with real bitcoin.
        div.panel-footer
          button.start-node START NODE »

</template>

<script>

</script>
