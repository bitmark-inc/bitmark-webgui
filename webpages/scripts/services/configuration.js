
app.service('configuration', function () {
    var configuration = {
        chain: "testing",
        bitmarkCliConfigFile: "",
        mineFee: 205000 // bitmark mine: 20000 bitcoin mine: 5000
    };

    return {
        getConfiguration: function () {
            return configuration;
        },
        setChain: function(chain) {
            configuration.chain = chain;
        },
        setBitmarkCliConfigFile: function(file){
            configuration.bitmarkCliConfigFile = file;
        }
    };
});
