package gurunavi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	APIEndpoint = "https://api.gnavi.co.jp/RestSearchAPI/v3/"
)

type GurunaviResponseBody struct {
	Attributes    *Attributes `json:"@attributes"`
	TotalHitCount int         `json:"total_hit_count"`
	HitPerPage    int         `json:"hit_per_page"`
	PageOffset    int         `json:"page_offset"`
	Rest          []*Rest     `json:"rest"`
	Error         []*Error    `json:"error"`
}
type Attributes struct {
	APIVersion string `json:"api_version"`
}
type Rest struct {
	RestAttributes *RestAttributes `json:"@attributes"`
	ID             string          `json:"id"`
	UpdateDate     time.Time       `json:"update_date"`
	Name           string          `json:"name"`
	NameKana       string          `json:"name_kana"`
	Latitude       string          `json:"latitude"`
	Longitude      string          `json:"longitude"`
	Category       string          `json:"category"`
	URL            string          `json:"url"`
	URLMobile      string          `json:"url_mobile"`
	CouponURL      *CouponURL      `json:"coupon_url"`
	ImageURL       *ImageURL       `json:"image_url"`
	Address        string          `json:"address"`
	Tel            string          `json:"tel"`
	TelSub         string          `json:"tel_sub"`
	Fax            string          `json:"fax"`
	Opentime       string          `json:"opentime"`
	Holiday        string          `json:"holiday"`
	Access         *Access         `json:"access"`
	ParkingLots    string          `json:"parking_lots"`
	Pr             *Pr             `json:"pr"`
	Code           *Code           `json:"code"`
	Budget         interface{}     `json:"budget"`
	Party          interface{}     `json:"party"`
	Lunch          interface{}     `json:"lunch"`
	CreditCard     string          `json:"credit_card"`
	EMoney         string          `json:"e_money"`
	Flags          *Flags          `json:"flags"`
}
type RestAttributes struct {
	Order int `json:"order"`
}
type CouponURL struct {
	Pc     string `json:"pc"`
	Mobile string `json:"mobile"`
}
type ImageURL struct {
	ShopImage1 string `json:"shop_image1"`
	ShopImage2 string `json:"shop_image2"`
	Qrcode     string `json:"qrcode"`
}
type Access struct {
	Line        string `json:"line"`
	Station     string `json:"station"`
	StationExit string `json:"station_exit"`
	Walk        string `json:"walk"`
	Note        string `json:"note"`
}
type Pr struct {
	PrShort string `json:"pr_short"`
	PrLong  string `json:"pr_long"`
}
type Code struct {
	Areacode      string   `json:"areacode"`
	Areaname      string   `json:"areaname"`
	Prefcode      string   `json:"prefcode"`
	Prefname      string   `json:"prefname"`
	AreacodeS     string   `json:"areacode_s"`
	AreanameS     string   `json:"areaname_s"`
	CategoryCodeL []string `json:"category_code_l"`
	CategoryNameL []string `json:"category_name_l"`
	CategoryCodeS []string `json:"category_code_s"`
	CategoryNameS []string `json:"category_name_s"`
}
type Flags struct {
	MobileSite   int `json:"mobile_site"`
	MobileCoupon int `json:"mobile_coupon"`
	PcCoupon     int `json:"pc_coupon"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SearchRestaurants(w string) (*GurunaviResponseBody, error) {
	v := url.Values{}
	v.Add("keyid", os.Getenv("GURUNAVI_ACCESS_KEY"))
	v.Add("freeword", w)

	// get で　request　を送る
	resp, err := http.Get(APIEndpoint + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	// TCPコネクション切断処理
	defer resp.Body.Close()

	// レスポンスを最後まで読み込む
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var g *GurunaviResponseBody

	if err := json.Unmarshal(body, &g); err != nil {
		return nil, err
	}

	return g, err
}
