package aftersales

// SCHEDULE_TPL 指派了师傅之后推送模板消息给师傅
const SCHEDULE_TPL = `
{
  "touser": "{{.OpenID}}",
  "template_id": "RvIkjIQOw7mf6oBCbMVhBQgE5cXNwoCQWy1L6w3gJhU",
  "url": "{{.URL}}",
  "topcolor": "#FF0000",
  "data": {
    "first": {
      "value": "您有新的订单，请尽快和用户约定具体上门服务时间并反馈",
      "color": "#173177"
    },
    "keyword1": {
      "value": "{{.ID}}",
      "color": "#173177"
    },
    "keyword2": {
      "value": "修水家电售后",
      "color": "#173177"
    },
    "keyword3": {
      "value": "{{.Now}}",
      "color": "#173177"
    },
    "remark": {
      "value": "请尽快处理， 为了保持良好的客户体验， 如果您超时未响应该笔订单可能会被分配给其他人了， 谢谢!",
      "color": "#173177"
    }
  }
}
`

// AUDIT_OK_TPL 审核成功通知模板
const AUDIT_OK_TPL = `
{
  "touser": "{{.OpenID}}",
  "template_id": "CoApRwtt1cD9o5aatsvf2yswTzGz4wHMb2uWFnyvWoY",
  "url": "{{.URL}}",
  "topcolor": "#FF0000",
  "data": {
    "first": {
      "value": "你好，订单{{.ID}}完成证明已经审核了",
      "color": "#173177"
    },
    "keyword1": {
      "value": "审核通过",
      "color": "#173177"
    },
    "keyword2": {
      "value": "{{.Now}}",
      "color": "#173177"
    },
    "remark": {
      "value": "你好订单{{.ID}}完成证明已经通过审核了， 按照约定， 该订单会进入冻结期，服务费用会在解冻之后可以提现， 如有疑问， 请咨询客服!",
      "color": "#173177"
    }
  }
}
`

// AUDIT_FAILED_TPL 审核不通过通知模板
const AUDIT_FAILED_TPL = `
{
  "touser": "{{.OpenID}}",
  "template_id": "CoApRwtt1cD9o5aatsvf2yswTzGz4wHMb2uWFnyvWoY",
  "url": "{{.URL}}",
  "topcolor": "#FF0000",
  "data": {
    "first": {
      "value": "你好，订单{{.ID}}完成证明审核不通过",
      "color": "#173177"
    },
    "keyword1": {
      "value": "审核不通过",
      "color": "#173177"
    },
    "keyword2": {
      "value": "{{.Now}}",
      "color": "#173177"
    },
    "remark": {
      "value": "你好订单{{.ID}}完成证明未通过审核， 如有疑问， 请咨询客服!",
      "color": "#173177"
    }
  }
}
`

// UNFREEZE_TPL 解冻通知
const UNFREEZE_TPL = `
{
  "touser": "{{.OpenID}}",
  "template_id": "CoApRwtt1cD9o5aatsvf2yswTzGz4wHMb2uWFnyvWoY",
  "url": "{{.URL}}",
  "topcolor": "#FF0000",
  "data": {
    "first": {
      "value": "你好，订单{{.ID}}已经解冻",
      "color": "#173177"
    },
    "keyword1": {
      "value": "解除冻结，服务费用可结算",
      "color": "#173177"
    },
    "keyword2": {
      "value": "{{.Now}}",
      "color": "#173177"
    },
    "remark": {
      "value": "你好订单{{.ID}}已经解冻， 该订单的服务费用转为可结算状态!",
      "color": "#173177"
    }
  }
}
`

type ModelForSchedule struct {
	OpenID string
	ID     string
	URL    string
	Now    string
}

type ModelForAudit struct {
	OpenID string
	ID     string
	URL    string
	Now    string
}

type Token struct {
	Token string `json:"token"`
}
