var Promise  = require('bluebird');
var OrderAPI = require('../api/OrderAPI');
var UserAPI  = require('../api/UserAPI');

module.exports = function(app) {
	// 登录验证中间件
	function requireLogin(req, res, next) {
		if (req.session && req.session.user) {
			return next();
		} else {
			res.redirect('/login');
		}
	}
	// 添加根路径路由 - 将用户重定向到 /admin
	app.get('/', function(req, res) {
		res.redirect('/admin');
	});
	// 显示登录页面
	app.get('/login', function(req, res) {
		res.render('admin/login', {
			layout: 'login' // 可选，如果您有特定的登录布局
		});
	});
	// 处理登录请求
	app.post('/login', function(req, res) {
		var username = req.body.username;
		var password = req.body.password;

		// 简单的输入验证
		if (!username || !password) {
			return res.render('admin/login', {
				error: '用户名和密码不能为空'
			});
		}

		UserAPI.authenticate(req, username, password)
			.then(function(user) {
				// 登录成功，保存用户信息到 session
				req.session.user = user;
				// 重定向到管理页面
				res.redirect('/admin');
			})
			.catch(function(err) {
				// 登录失败，显示错误信息
				res.render('admin/login', {
					error: err.message || '用户名或密码错误'
				});
			});
	});

	app.get('/admin', requireLogin, function(req, res) {
		Promise.all([
			OrderAPI.getTodayOrderCount(req),
			OrderAPI.getTodaySale(req),
			OrderAPI.getTotalOrderCount(req),
			OrderAPI.getTotalSale(req)
		])
		.then(function(arr) {
			res.locals.data = {
				todayOrderCount : arr[0].count,
			    todayTotalSale  : arr[1].sale,
			    totalOrderCount : arr[2].count,
			    totalSale       : arr[3].sale,
				user            : req.session.user  // 传递用户信息到前端

			};
			res.render('admin/index');
		})
		.catch(function(err) {
			console.log(err)
	        res.status(500);
		    res.render('error', {
		        message : err.message,
		        error   : {}
		    });
		});
	});
	// 添加登出功能
	app.get('/logout', function(req, res) {
		req.session.destroy();
		res.redirect('/login');
	});

};
