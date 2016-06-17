
app.service('configuration', function () {
    var configuration = {
        chain: "testing",
        bitmarkCliConfigFile: ""
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
