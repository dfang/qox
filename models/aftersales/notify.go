package aftersales

var tpl1 = `
{
  "touser": "{{.OpenID}}",
  "template_id": "RvIkjIQOw7mf6oBCbMVhBQgE5cXNwoCQWy1L6w3gJhU",
  "url": "{{.URL}}",
  "topcolor": "#FF0000",
  "data": {
    "first": {
      "value": "您有新的工单，请尽快和用户约定具体上门服务时间并反馈",
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
      "value": "{{.Date}}",
      "color": "#173177"
    },
    "remark": {
      "value": "请尽快处理",
      "color": "#173177"
    }
  }
}
`

type ModelForTpl1 struct {
	OpenID string
	ID     string
	URL    string
	Date   string
}

type Token struct {
	Token string `json:"token"`
}
