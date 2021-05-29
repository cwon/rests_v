package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/buger/jsonparser"
	"github.com/gen2brain/beeep"
	"github.com/gocolly/colly"
)

//var text_payload = `[{"operationName":"vaccineList","variables":{"input":{"keyword":"코로나백신위탁의료기관","x":"127.1399977","y":"37.264899"},"businessesInput":{"start":0,"display":100,"deviceType":"mobile","x":"127.1399977","y":"37.264899","bounds":"127.1017601;37.2432423;127.1782353;37.2865495","sortingOrder":"distance"},"isNmap":false,"isBounds":false},"query":"query vaccineList($input: RestsInput, $businessesInput: RestsBusinessesInput, $isNmap: Boolean!, $isBounds: Boolean!) {\n  rests(input: $input) {\n    businesses(input: $businessesInput) {\n      total\n      vaccineLastSave\n      isUpdateDelayed\n      items {\n        id\n        name\n        dbType\n        phone\n        virtualPhone\n        hasBooking\n        hasNPay\n        bookingReviewCount\n        description\n        distance\n        commonAddress\n        roadAddress\n        address\n        imageUrl\n        imageCount\n        tags\n        distance\n        promotionTitle\n        category\n        routeUrl\n        businessHours\n        x\n        y\n        imageMarker @include(if: $isNmap) {\n          marker\n          markerSelected\n          __typename\n        }\n        markerLabel @include(if: $isNmap) {\n          text\n          style\n          __typename\n        }\n        isDelivery\n        isTakeOut\n        isPreOrder\n        isTableOrder\n        naverBookingCategory\n        bookingDisplayName\n        bookingBusinessId\n        bookingVisitId\n        bookingPickupId\n        vaccineQuantity {\n          quantity\n          quantityStatus\n          vaccineType\n          vaccineOrganizationCode\n          __typename\n        }\n        __typename\n      }\n      optionsForMap @include(if: $isBounds) {\n        maxZoom\n        minZoom\n        includeMyLocation\n        maxIncludePoiCount\n        center\n        __typename\n      }\n      __typename\n    }\n    queryResult {\n      keyword\n      vaccineFilter\n      categories\n      region\n      isBrandList\n      filterBooking\n      hasNearQuery\n      isPublicMask\n      __typename\n    }\n    __typename\n  }\n}\n"}]`
var text_payload = `[{"operationName":"vaccineList","variables":{"input":{"keyword":"코로나백신위탁의료기관","x":"Longitude","y":"Latitude"},"businessesInput":{"start":0,"display":100,"deviceType":"mobile","x":"Longitude","y":"Latitude","bounds":"LonB1;LatB1;LonB2;LatB2","sortingOrder":"distance"},"isNmap":false,"isBounds":false},"query":"query vaccineList($input: RestsInput, $businessesInput: RestsBusinessesInput, $isNmap: Boolean!, $isBounds: Boolean!) {\n  rests(input: $input) {\n    businesses(input: $businessesInput) {\n      total\n      vaccineLastSave\n      isUpdateDelayed\n      items {\n        id\n        name\n        dbType\n        phone\n        virtualPhone\n        hasBooking\n        hasNPay\n        bookingReviewCount\n        description\n        distance\n        commonAddress\n        roadAddress\n        address\n        imageUrl\n        imageCount\n        tags\n        distance\n        promotionTitle\n        category\n        routeUrl\n        businessHours\n        x\n        y\n        imageMarker @include(if: $isNmap) {\n          marker\n          markerSelected\n          __typename\n        }\n        markerLabel @include(if: $isNmap) {\n          text\n          style\n          __typename\n        }\n        isDelivery\n        isTakeOut\n        isPreOrder\n        isTableOrder\n        naverBookingCategory\n        bookingDisplayName\n        bookingBusinessId\n        bookingVisitId\n        bookingPickupId\n        vaccineQuantity {\n          quantity\n          quantityStatus\n          vaccineType\n          vaccineOrganizationCode\n          __typename\n        }\n        __typename\n      }\n      optionsForMap @include(if: $isBounds) {\n        maxZoom\n        minZoom\n        includeMyLocation\n        maxIncludePoiCount\n        center\n        __typename\n      }\n      __typename\n    }\n    queryResult {\n      keyword\n      vaccineFilter\n      categories\n      region\n      isBrandList\n      filterBooking\n      hasNearQuery\n      isPublicMask\n      __typename\n    }\n    __typename\n  }\n}\n"}]`

func main() {
	ch := make(chan bool)

	a := app.New()
	w := a.NewWindow("Hello")
	w.Resize(fyne.NewSize(400, 200))
	w.SetFixedSize(true)

	hello := widget.NewLabel("Rest Vaccine Monitoring...")
	entry1 := widget.NewEntry()
	entry1.SetPlaceHolder("Latitude - approximately, 37......")
	entry2 := widget.NewEntry()
	entry2.SetPlaceHolder("Longitude - approximately, 127....")
	entry3 := widget.NewEntry()
	entry3.SetPlaceHolder("Interval second")
	startBtn := widget.NewButton("Start", func() {})
	startBtn.OnTapped = func() {
		if startBtn.Text == "Start" {
			entry1.Disable()
			entry2.Disable()
			entry3.Disable()
			startBtn.SetText("Stop")
			sec, err := strconv.Atoi(entry3.Text)
			latitude, err := strconv.ParseFloat(entry1.Text, 32)
			longitude, err := strconv.ParseFloat(entry2.Text, 32)
			if err == nil {
				go get_rest_vaccine_data(ch, uint(sec), latitude, longitude)
			}
		} else {
			ch <- false
			entry1.Enable()
			entry2.Enable()
			entry3.Enable()
			startBtn.SetText("Start")
		}
	}

	w.SetContent(container.NewVBox(
		hello,
		entry1,
		entry2,
		entry3,
		startBtn,
	))

	w.ShowAndRun()
}

func deepCopy(s string) string {
	b := make([]byte, len(s))
	copy(b, s)
	return *(*string)(unsafe.Pointer(&b))
}

func get_rest_vaccine_data(quit chan bool, interval_time uint, latitude float64, longitude float64) {
	ticker := time.NewTicker(time.Duration(interval_time) * time.Second)
	payload := deepCopy(text_payload)

	lon := fmt.Sprintf("%f", longitude)
	lat := fmt.Sprintf("%f", latitude)

	b1_lon := fmt.Sprintf("%f", longitude-0.038)
	b1_lat := fmt.Sprintf("%f", latitude-0.021)

	b2_lon := fmt.Sprintf("%f", longitude+0.038)
	b2_lat := fmt.Sprintf("%f", latitude+0.021)

	payload = strings.ReplaceAll(payload, "Longitude", lon)
	payload = strings.ReplaceAll(payload, "Latitude", lat)
	payload = strings.ReplaceAll(payload, "LonB1", b1_lon)
	payload = strings.ReplaceAll(payload, "LatB1", b1_lat)
	payload = strings.ReplaceAll(payload, "LonB2", b2_lon)
	payload = strings.ReplaceAll(payload, "LatB2", b2_lat)

	for {
		select {
		case <-ticker.C:
			c := colly.NewCollector()
			c.OnRequest(func(r *colly.Request) {
				r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
			})
			c.OnResponse(func(r *colly.Response) {
				jsonparser.ArrayEach(r.Body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					name, _ := jsonparser.GetString(value, "name")
					road, _ := jsonparser.GetString(value, "roadAddress")
					vq, err := jsonparser.GetString(value, "vaccineQuantity", "quantity")
					q, err := strconv.Atoi(vq)
					if err == nil {
						if q > 0 {
							beeep.Notify("vaccineQuantity > 0 !", name+road, "")
						}
					}

				}, "[0]", "data", "rests", "businesses", "items")
				fmt.Println("iter")
			})
			c.PostRaw("https://api.place.naver.com/graphql", []byte(payload))
		case <-quit:
			fmt.Println("goodbye")
			ticker.Stop()
			return
		}
	}
}
