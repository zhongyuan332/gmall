var request   = require('../net/request');
var Promise   = require('bluebird');
var config    = require('../config');
var ErrorCode = require('../model/ErrorCode');

/*
 * 用户认证
 */
function authenticate(req, username, password) {
    return new Promise(function(resolve, reject) {
        // 这里假设后端已经提供了认证 API
        var url = config.api.userAuth;
        request.postJSON({
            client: req,
            uri: url,
            body: {
                username: username,
                password: password
            }
        }, function(error, response, data) {
            console.log('UserAPI.js authenticate', error, data);
            if (data && data.errNo != ErrorCode.SUCCESS) {
                reject(data);
            } else if (!error && response.statusCode == 200 && data && data.errNo === ErrorCode.SUCCESS) {
                resolve(data.data);
            } else {
                reject(error);
            }
        });
    });
}

module.exports = {
    authenticate: authenticate
};