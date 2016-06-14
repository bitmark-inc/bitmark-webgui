app.service('configuration', function () {
    var configuration = {
        chain: "testing"
    };

    return {
        getConfiguration: function () {
            return configuration;
        },
        setChain: function(chain) {
            configuration.chain = chain;
        }
    };
});
