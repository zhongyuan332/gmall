var Promise  = require('bluebird');
var OrderAPI = require('../api/OrderAPI');

module.exports = function(app) {
	// 添加根路径路由 - 将用户重定向到 /admin
	app.get('/', function(req, res) {
		res.redirect('/admin');
	});

	app.get('/admin', function(req, res) {
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
			    totalSale       : arr[3].sale
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
};
