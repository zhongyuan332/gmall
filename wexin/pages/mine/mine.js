var login = require('../../common/login.js');

Page({
  data: {
    userInfo: null
  },
  onUserInfoCallback(userInfo) {
    this.setData({
      userInfo: userInfo
    })
  },
  onLoad: function () {
    console.log("mine page onLoad start")
    var app = getApp();
    var userInfo = app.globalData.userInfo;
    console.log("mine page, userInfo:", userInfo)
    if (!userInfo) {
      login.addLoginResponder(this.onUserInfoCallback.bind(this));
    } else {
      this.setData({
        userInfo: userInfo
      })
    }
  }
})