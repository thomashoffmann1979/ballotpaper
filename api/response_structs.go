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



type StimmzettelResponse struct {
	Msg                 string `json:"msg"`
	Success             bool   `json:"success"`
	Errors              []any  `json:"errors"`
	Warnings            []any  `json:"warnings"`
	DsxRestAPIGetResult string `json:"dsx_rest_api_get_result"`
	Data                []struct {
		TableName                             string `json:"__table_name"`
		ID                                    string `json:"__id"`
		Displayfield                          string `json:"__displayfield"`
		Ridx                                  string `json:"ridx"`
		ID0                                   string `json:"id"`
		Name                                  string `json:"name"`
		Aktiv                                 string `json:"aktiv"`
		InsertDate                            string `json:"insert_date"`
		InsertTime                            string `json:"insert_time"`
		UpdateDate                            string `json:"update_date"`
		UpdateTime                            string `json:"update_time"`
		Login                                 string `json:"login"`
		Wahlgruppe                            string `json:"wahlgruppe"`
		Wahlbezirk                            string `json:"wahlbezirk"`
		Wahltyp                               string `json:"wahltyp"`
		Sitze                                 string `json:"sitze"`
		Anzahl10                              string `json:"anzahl_10"`
		Zaehlung1                             string `json:"zaehlung_1"`
		Zaehlung2                             string `json:"zaehlung_2"`
		Zaehlung3                             string `json:"zaehlung_3"`
		Zaehlung4                             string `json:"zaehlung_4"`
		Zaehlung5                             string `json:"zaehlung_5"`
		Zaehlung6                             string `json:"zaehlung_6"`
		Zaehlung7                             string `json:"zaehlung_7"`
		Zaehlung8                             string `json:"zaehlung_8"`
		Zaehlung9                             string `json:"zaehlung_9"`
		Sitzbindung                           string `json:"sitzbindung"`
		LaufendeNummer124NachZuordnungWgWb    any    `json:"laufende_nummer_1_24__nach_zuordnung_wg_wb"`
		AnzahlNotwendigeBewerberJeKombinummer any    `json:"anzahl_notwendige_bewerber_je_kombinummer"`
		DsCount                               string `json:"ds_count"`
		Ungueltig                             any    `json:"ungueltig"`
		Farbe                                 string `json:"farbe"`
		Typtitel                              any    `json:"typtitel"`
		KandidatTextEinzel                    any    `json:"kandidat_text_einzel"`
		KandidatTextMehr                      any    `json:"kandidat_text_mehr"`
		Bh1                                   any    `json:"bh1"`
		Bh2                                   any    `json:"bh2"`
		Xlink                                 any    `json:"xlink"`
		Rownumber                             int    `json:"__rownumber"`
	} `json:"data"`
	Total int `json:"total"`
}


type KandidatenResponse struct {
	Msg                 string `json:"msg"`
	Success             bool   `json:"success"`
	Errors              []any  `json:"errors"`
	Warnings            []any  `json:"warnings"`
	DsxRestAPIGetResult string `json:"dsx_rest_api_get_result"`
	Data                []struct {
		TableName                             string `json:"__table_name"`
		ID                                    string `json:"__id"`
		Barcode                               string `json:"barcode"`
		Displayfield                          string `json:"__displayfield"`
	} `json:"data"`
	Total int `json:"total"`
}