var config = require("../config/config.js");

var login = {
  loginResponders: [],
  addLoginResponder: function (responder) {
    this.loginResponders.push(responder);
  },
  login: function (app) {
    var self = this;
    var resData = {};
    var jsCodeDone, userInfoDone;

    function setUserInfo() {
      wx.request({
        url: config.api.setWeAppUser,
        data: {
          encryptedData: resData.encryptedData,
          iv: resData.iv
        },
        header: {
          'content-type': 'application/json',
          'Cookie': "sid=" + resData.sid
        },
        method: "POST",
        success: function (res) {
          app.globalData.userInfo = resData.userInfo;
          app.globalData.encryptedData = resData.encryptedData;
          app.globalData.iv = resData.iv;
          app.globalData.sid = resData.sid;
          for (var i = 0; i < self.loginResponders.length; i++) {
            self.loginResponders[i]();
          }
        }
      });
    }

    wx.login({
      success: function (res) {
        console.log("微信登录成功，返回:", res)
        if (res.code) {
          wx.request({
            url: config.api.weAppLogin,
            data: {
              code: res.code
            },
            success: function (res) {
              console.log("登录成功，设置后台session成功,", res)
              resData.sid = res.data.data.sid;
              jsCodeDone = true;
              jsCodeDone && userInfoDone && setUserInfo();
            }
          });

          wx.getUserInfo({
            success: function (res) {
              resData.userInfo = res.userInfo;
              resData.encryptedData = res.encryptedData;
              resData.iv = res.iv;
              userInfoDone = true;
              jsCodeDone && userInfoDone && setUserInfo();
            },
            fail: function (data) {
              console.log(data);
            }
          });
        }
      }
    });
  },

  logout: function () {

  }
}

module.exports = login;