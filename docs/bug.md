# README


1. [scope with default:true 不能用](https://github.com/qor/admin/issues/199)
  设置了default:true 的scope，会影响到action with user input的更新操作
  禁用 default: true 的scope 才能在action 里成功 更新记录
