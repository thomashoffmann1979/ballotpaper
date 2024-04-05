package api


type LoginResponse struct {
	Msg      string `json:"msg"`
	Success  bool   `json:"success"`
	Errors   []any  `json:"errors"`
	Warnings []any  `json:"warnings"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Client   string `json:"client"`
	Clients  []struct {
		Client string `json:"client"`
	} `json:"clients"`
	Dbaccess bool `json:"dbaccess"`
}


type PingResponse struct {
	Msg      string `json:"msg"`
	Success  bool   `json:"success"`
	Errors   []any  `json:"errors"`
	Warnings []any  `json:"warnings"`
	Username string `json:"username"`
	Clients  []struct {
		Client string `json:"client"`
	} `json:"clients"`
	Client       string `json:"client"`
	Fullname     string `json:"fullname"`
	Gst          string `json:"gst"`
	Bkr          string `json:"bkr"`
	Gstavatar    string `json:"gstavatar"`
	Bkravatar    string `json:"bkravatar"`
	Avatar       string `json:"avatar"`
	Clientavatar string `json:"clientavatar"`
}