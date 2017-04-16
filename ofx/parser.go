package ofx

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type OFX struct {
	XMLName      xml.Name      `xml:"OFX"`
	Balance      float64       `xml:"BANKMSGSRSV1>STMTTRNRS>STMTRS>AVAILBAL>BALAMT"`
	Transactions []Transaction `xml:"BANKMSGSRSV1>STMTTRNRS>STMTRS>BANKTRANLIST>STMTTRN"`
}

type Transaction struct {
	XMLName xml.Name `xml:"STMTTRN"`

	Number      int
	TxType      string `xml:"TRNTYPE"`
	Date        string `xml:"DTPOSTED"`
	Amount      string `xml:"TRNAMT"`
	ID          string `xml:"FITID"`
	Description string `xml:"NAME"`
}

func (t *Transaction) GetAmount() float64 {
	if t.Amount == "" {
		return 0.0
	} else {

		str := strings.Replace(t.Amount, ",", ".", -1)
		re := regexp.MustCompile("(\\d)\\s(\\d)")
		str = re.ReplaceAllString(str, "${1}${2}")
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return 0
		}
		return v
	}
}

func (t *Transaction) GetDate() time.Time {
	v, err := time.Parse("20060102", t.Date)
	if err != nil {
		return time.Time{}
	}
	return v
}

func Parse(filepath string) (*OFX, error) {
	xmlFile, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer xmlFile.Close()
	var o OFX
	b, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading:", err)
		return nil, err
	}
	xml.Unmarshal(b, &o)
	return &o, err
}
