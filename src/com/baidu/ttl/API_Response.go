package ttl

type API_Response struct {
	Err_no  int        `json:"err_no"`
	Err_msg string     `json:"err_msg"`
	Sn      string     `json:"sn"`
	Idx     int        `json:"idx"`
}