var config = require('./config/config.js');
var login  = require('./common/login.js');

App({
    onLaunch: function() {
      console.log("初始化登录")
        login.login(this);
    },
    globalData: {
        userInfo: null,
        encryptedData: "",
        iv: "",
        sid: ""
    }
})