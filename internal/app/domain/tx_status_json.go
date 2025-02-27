package domain

import "encoding/xml"

type TxStatusRabbitJson struct {
	TxID   int64  `json:"tx_id"`
	Status string `json:"status"`
}

type TxStatusXML struct {
	XMLName xml.Name `xml:"tx_status"`
	TxID    int64    `xml:"tx_id"`
	Status  string   `xml:"status"`
}
