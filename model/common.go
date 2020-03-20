package model

//删除
type DeleteReqJson struct {
	Ids []int `json:"ids"` //id数组
	Id  int   `json:"id"`  //id数组
}
