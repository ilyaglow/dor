package dor

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	yandexRadarTop10k = "https://radar.yandex.ru/export?options=%7B%22title%22%3A%22%D0%AF%D0%BD%D0%B4%D0%B5%D0%BA%D1%81.%D0%A0%D0%B0%D0%B4%D0%B0%D1%80%2B-%2B%D0%A2%D0%BE%D0%BF%2B%D0%98%D0%BD%D1%82%D0%B5%D1%80%D0%BD%D0%B5%D1%82-%D0%BF%D1%80%D0%BE%D0%B5%D0%BA%D1%82%D0%BE%D0%B2%22%2C%22exportFormat%22%3A%22csv%22%2C%22limit%22%3A0%2C%22offset%22%3A1000000%7D&mode=top-sites"
)

// YandexRadarIngester represents Ingester implementation for Yandex Radar.
type YandexRadarIngester struct {
	IngesterConf
}

// NewYandexRadar bootstraps YandexRadarIngester.
func NewYandexRadar() *YandexRadarIngester {
	return &YandexRadarIngester{
		IngesterConf: IngesterConf{
			Description: "yandex-radar",
		},
	}

}

// Do implements Ingester Do func with the data.
func (in *YandexRadarIngester) Do() (chan *Entry, error) {
	in.Timestamp = time.Now().UTC()
	ch := make(chan *Entry)

	go func() {
		defer close(ch)
		resp, err := http.Get(yandexRadarTop10k)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("%s: %s downloaded successfully", in.Description, yandexRadarTop10k)
		defer resp.Body.Close()

		r := csv.NewReader(resp.Body)
		r.LazyQuotes = true
		// read the header
		// "Название ресурса","URL-адрес ресурса","Тематики ресурса","Тип ресурса","Медиахолдинг","Данные Метрики","Посетители (кросс-девайс)","Посетители (браузер)","Среднее время","Доля пользователей приложения","Дневная аудитория"
		_, err = r.Read()
		if err != nil {
			log.Println(err)
			return
		}
		var i uint32
		now := time.Now()
		for {
			i++
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Println(err)
				return
			}

			ch <- &Entry{
				Domain:  record[1],
				Rank:    i,
				Date:    now,
				Source:  in.Description,
				RawData: strings.Join(record, ","),
			}
		}
	}()

	return ch, nil
}
