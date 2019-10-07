moment.locale('en', {
  dow: 1, // Monday is the first day of the week.
})
moment.locale('zh-cn', {
  dow: 1, // Monday is the first day of the week.
})

var OrderChart, UsersChart;
function RenderChart(ordersData, usersData) {
  Chart.defaults.global.responsive = true;

  var orderDateLables = [];
  var orderCounts = [];
  for (var i = 0; i < ordersData.length; i++) {
    orderDateLables.push(ordersData[i].Date.substring(5, 10));
    orderCounts.push(ordersData[i].Total)
  }
  if (OrderChart) {
    OrderChart.destroy();
  }
  var orders_context = document.getElementById("orders_report").getContext("2d");
  var orders_data = ChartData(orderDateLables, orderCounts);
  OrderChart = new Chart(orders_context).Line(orders_data, "");

  // var usersDateLables = [];
  // var usersCounts = [];
  // for (var i = 0; i < usersData.length; i++) {
  //   usersDateLables.push(usersData[i].Date.substring(5, 10));
  //   usersCounts.push(usersData[i].Total)
  // }
  // if (UsersChart) {
  //   UsersChart.destroy();
  // }
  // var users_context = document.getElementById("users_report").getContext("2d");
  // var users_data = ChartData(usersDateLables, usersCounts);
  // UsersChart = new Chart(users_context).Bar(users_data, "");
}

function ChartData(lables, counts) {
  var chartData = {
    labels: lables,
    datasets: [
      {
        label: "Users Report",
        fillColor: "rgba(151,187,205,0.2)",
        strokeColor: "rgba(151,187,205,1)",
        pointColor: "rgba(151,187,205,1)",
        pointStrokeColor: "#fff",
        pointHighlightFill: "#fff",
        pointHighlightStroke: "rgba(151,187,205,1)",
        data: counts
      }
    ]
  };
  return chartData;
}

Date.prototype.Format = function (fmt) {
  var o = {
    "M+": this.getMonth() + 1,
    "d+": this.getDate(),
    "h+": this.getHours(),
    "m+": this.getMinutes(),
    "s+": this.getSeconds(),
    "q+": Math.floor((this.getMonth() + 3) / 3),
    "S": this.getMilliseconds()
  };
  if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
  for (var k in o)
    if (new RegExp("(" + k + ")").test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
  return fmt;
}

Date.prototype.AddDate = function (add) {
  var date = this.valueOf();
  date = date + add * 24 * 60 * 60 * 1000
  date = new Date(date)
  return date;
}

// 给body标签加一个class 以方便区别是什么模块什么页面
function addClassToBody() {
  // 比如页面 admin/orders/1
  path = window.location.pathname.split("/")[2] || 'dashboard'
  className = 'qor-' + path
  $('body').addClass(className)
}

// qor dashboard
$(document).ready(function () {
  addClassToBody()

  // 订单列表页 禁止slideout 打开订单详情
  if ($('.qor-orders .qor-table-container table').length > 0) {
    $(document).off("click.qor.openUrl", "[data-url]")
  }

  $("#startDate").val(moment().subtract(10, 'days').format('YYYY-MM-DD'));
  $("#endDate").val(moment().format('YYYY-MM-DD'));
  $(".j-update-record").click(function () {
    $.getJSON("/admin/reports.json", { startDate: $("#startDate").val(), endDate: $("#endDate").val() }, function (jsonData) {
      RenderChart(jsonData.Orders, jsonData.Users);
    });
  });
  $(".j-update-record").click();

  $(".this-week-reports").click(function () {
    var start = moment().startOf('week').format('YYYY-MM-DD')
    var end = moment().format('YYYY-MM-DD')
    $("#startDate").val(start);
    $("#endDate").val(end);
    $(".j-update-record").click();
    $(this).blur();
  });

  $(".last-week-reports").click(function () {
    var end = moment().startOf('week').subtract(1, 'days').format('YYYY-MM-DD')
    var start = moment().startOf('week').subtract(1, 'days').startOf('week').format('YYYY-MM-DD')
    $("#startDate").val(start);
    $("#endDate").val(end);
    $(".j-update-record").click();
    $(this).blur();
  });

  $(".this-month-reports").click(function () {
    var start = moment().startOf('month').format('YYYY-MM-DD')
    var end = moment().format('YYYY-MM-DD')
    $("#startDate").val(start);
    $("#endDate").val(end);
    $(".j-update-record").click();
    $(this).blur();
  });

  $(".last-month-reports").click(function () {
    var start = moment().subtract(1, 'months').startOf('month').format('YYYY-MM-DD');
    var end = moment().subtract(1, 'months').endOf('month').format('YYYY-MM-DD');
    $("#startDate").val(start);
    $("#endDate").val(end);
    $(".j-update-record").click();
    $(this).blur();
  });

  $(".this-year-reports").click(function () {
    var start = moment().startOf('year').format('YYYY-MM-DD');
    var end = moment().format('YYYY-MM-DD');
    $("#startDate").val(start);
    $("#endDate").val(end);
    $(".j-update-record").click();
    $(this).blur();
  });

  $(".this-year2-reports").click(function () {
    var start = moment().startOf('year').format('YYYY-MM');
    var end = moment().format('YYYY-MM');
    $("#startDate").val(start);
    $("#endDate").val(end);
    $(".j-update-record").click();
    $(this).blur();
  });

  // $('#startDate').on('onOk', function () {
  //   console.log("fuck")
  // })

  $('.qor-datepicker').find('.qor-datepicker__save').on('click', function (e) {
    console.log('fuck');
  });

  // document.getElementById('startDate').addEventListener('onOk', function () {
  //   console.log("fuck")
  // });

  // document.getElementById('endDate').addEventListener('onOk', function () {
  //   this.value = x.time.toString();
  // });
});



