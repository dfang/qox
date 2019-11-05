// qor admin 是如何打包的，需要研究下这个
// https://github.com/qor/qor/blob/master/gulpfile.js

// 暂时需要引入https://github.com/ccampbell/mousetrap/blob/master/mousetrap.min.js
// 本应该打包到vendors.js, 暂时直接放这里

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


  if ($('.qor-balances').length > 0) {
    // hide big new button
    $('.qor-button--new').hide()

    // disable click to open url in tables
    $(document).off("click.qor.openUrl", "[data-url]")

    // hide qor-actions dropdown
    $('.qor-table__actions').hide()

    // hide search button
    $('.qor-search__label').hide()
  }


  if ($('.qor-settlements').length > 0) {
    // $(document).off("click.qor.openUrl", ".qor-table-container tr[data-url]")
    // $(document).off("click.qor.openUrl", "[data-url]")

    // hide search button
    $('.qor-search__label').hide()

    // hide qor-actions dropdown
    $('.qor-table__actions').hide()
  }

  if ($('.qor-brands, .qor-servcie_types, .qor-wechat_profiles').length > 0) {
    // disable click to open url
    $(document).off("click.qor.openUrl", "[data-url]")
  }
});


$(document).ready(function () {
  // restore Drawer state from cookie
  var x = getCookie('drawer_state');
  if (x == "1") {
    $('.mdl-layout').removeClass('hidden-drawer')
  } else {
    $('.mdl-layout').addClass('hidden-drawer')
  }

  $(document).on('click', '.mdl-layout__drawer-button', function (e) {
    showDrawer()
  });

  $(document).on('click', '.sidebar-footer, .sidebar-header a:last', function (e) {
    hideDrawer()
  });

  // go to dashboard
  Mousetrap.bind('g d', function () {
    window.location.href = "/admin"
  });

  // go to aftersales
  Mousetrap.bind('g a', function () {
    window.location.href = "/admin/aftersales"
  });

  Mousetrap.bind('g s', function () {
    window.location.href = "/admin/settlements"
  });

  Mousetrap.bind('g b', function () {
    window.location.href = "/admin/balances"
  });

  // New action
  Mousetrap.bind('n', function () {
    $('.qor-button--new').click()
  });

  Mousetrap.bind('s', function () {
    toggleDrawer()
  });

  Mousetrap.bind('d', function () {
    toggleDrawer()
  });

  // notifictions
  Mousetrap.bind('g n', function () {
    $('.qor-notifications__badges:last').click()
  });

  // help
  Mousetrap.bind('g h', function () {
    $('.qor-notifications__badges:first').click()
  });
})

function toggleDrawer() {
  var x = getCookie('drawer_state');
  if (x == "1") {
    hideDrawer()
  } else {
    showDrawer()
  }
}

function showDrawer() {
  $('.mdl-layout').removeClass('hidden-drawer')
  setCookie('drawer_state', '1', 365)
}

function hideDrawer() {
  console.log("hide drawer")
  $('.mdl-layout').addClass('hidden-drawer')
  $('.mdl-layout__obfuscator.is-visible').hide("slow")

  setCookie('drawer_state', '0', 365)
}

// setCookie('hide_drawer', '1', 365);
function setCookie(name, value, days) {
  var expires = "";
  if (days) {
    var date = new Date();
    date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
    expires = "; expires=" + date.toUTCString();
  }
  document.cookie = name + "=" + (value || "") + expires + "; path=/";
}

function getCookie(name) {
  var nameEQ = name + "=";
  var ca = document.cookie.split(';');
  for (var i = 0; i < ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == ' ') c = c.substring(1, c.length);
    if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length, c.length);
  }
  return null;
}
