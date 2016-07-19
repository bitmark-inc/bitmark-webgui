app.service('utils', function ($log) {
    var checkErrorObj = function(obj){
        var fields = ["show", "msg"];
        for(var i=0; i<fields.length; i++) {
            if(obj[fields[i]] == undefined){
                $log.error("error obj lack of field: "+fields[i]);
                return false;
            }
        }
        return true;
    };

    return {
        // obj format:
        // {
        //     show: true,
        //     msg: ""
        // }
        setErrorMsg: function (obj, show, msg) {
            if(!checkErrorObj(obj)){
                return;
            }
            obj.msg = msg;
            obj.show = show;
        }
    };
});
