package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"math"
	"strconv"

	"github.com/golang-geo/geo"
)

type kml struct {
	Document Document `xml:"Document"`
}
type Document struct {
	Name      string      `xml:"name"`
	Open      string      `xml:"open"`
	Style     []Style     `xml:"Style"`
	StyleMap  []StyleMap  `xml:"StyleMap"`
	Folder    []Folder    `xml:"Folder"`
	Placemark []Placemark `xml:"Placemark"`
}
type Style struct {
	Id        string    `xml:"id,attr"`
	IconStyle IconStyle `xml:"IconStyle"`
	LineStyle LineStyle `xml:"LineStyle"`
}
type IconStyle struct {
	Scale   string  `xml:"scale"`
	Icon    Icon    `xml:"Icon"`
	HotSpot HotSpot `xml:"hotSpot"`
}

type HotSpot struct {
	X      string `xml:"x,attr"`
	Y      string `xml:"y,attr"`
	Xunits string `xml:"xunits,attr"`
	Yunits string `xml:"yunits,attr"`
}
type Icon struct {
	Href string `xml:"href"`
}
type LineStyle struct {
	Color string `xml:"color"`
}
type Folder struct {
	Name        string      `xml:"name"`
	Description string      `xml:"description"`
	Placemark   []Placemark `xml:"Placemark"`
}
type Placemark struct {
	Name        string     `xml:"name"`
	Description string     `xml:"description"`
	LookAt      LookAt     `xml:"LookAt"`
	StyleUrl    string     `xml:"styleUrl"`
	LineString  LineString `xml:"LineString"`
	Point       Point      `xml:"Point"`
}
type LookAt struct {
	Longitude      string `xml:"longitude"`
	Latitude       string `xml:"latitude"`
	Altitude       string `xml:"altitude"`
	Heading        string `xml:"heading"`
	Tilt           string `xml:"tilt"`
	Range          string `xml:"range"`
	GxaltitudeMode string `xml:"http://www.google.com/kml/ext/2.2 altitudeMode"`
}
type Point struct {
	GxdrawOrder string `xml:"http://www.google.com/kml/ext/2.2 drawOrder"`
	Coordinates string `xml:"coordinates"`
}
type LineString struct {
	Tessellate  string `xml:"tessellate"`
	Coordinates string `xml:"coordinates"`
}
type StyleMap struct {
	Id   string `xml:"id,attr"`
	Pair []Pair `xml:"Pair"`
}
type Pair struct {
	Key      string `xml:"key"`
	StyleUrl string `xml:"styleUrl"`
}

func splitLink(s string) []string {

	s = strings.Replace(s, "0 ", "", -1)
	s = strings.TrimSpace(s)
	x := strings.Split(s, ",")
	return x[:len(x)-1]
}
func splitStr(s string) []string {
	s = strings.TrimSpace(s)
	x := strings.Split(s, ",")
	return x[:len(x)-1]
}
func FlushTestDB(s *geo.SQLMapper) {
	s.SqlConn.Exec("DELETE FROM points;")
}
func RoundFloat(x float64, prec int) float64 {
	frep := strconv.FormatFloat(x, 'g', prec, 64)
	f, _ := strconv.ParseFloat(frep, 64)
	return f
}
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
func NearestLinestringCo(s *geo.SQLMapper, longtitude1 float64, latitude1 float64) (float64, float64, float64) {

	RadiusData := []float64{}
	RadiusDistance := []float64{}
	var la, ln, d float64
	origin := geo.NewPoint(longtitude1, latitude1)
	res, err := s.PointsWithinRadius(origin, 0.4)
	RadiusData = RadiusData[:0]
	if err != nil {
		panic(err)
	}

	for res.Next() {
		err = res.Scan(&ln, &la)
		if err != nil {
			panic(err)
		}
		if ln != longtitude1 && la != latitude1 {
			RadiusData = append(RadiusData, ln)
			RadiusData = append(RadiusData, la)
		}
	}

	RadiusDistance = RadiusDistance[:0]
	for c := 0; c < len(RadiusData); c = c + 2 {
		lo1 := RadiusData[c]
		la1 := RadiusData[c+1]

		d = Distance(latitude1, longtitude1, la1, lo1)

		RadiusDistance = append(RadiusDistance, d)
		RadiusDistance = append(RadiusDistance, lo1)
		RadiusDistance = append(RadiusDistance, la1)

	}

	Small := RadiusDistance[0]
	long2 := RadiusDistance[1]
	lati2 := RadiusDistance[2]
	for k := 3; k < len(RadiusDistance); k = k + 3 {
		if Small > RadiusDistance[k] {
			Small = RadiusDistance[k]
			long2 = RadiusDistance[k+1]
			lati2 = RadiusDistance[k+2]

		}
	}

	return Small, long2, lati2
}
func Distance(lat1, lon1, lat2, lon2 float64) float64 {

	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
	// km:=m/1000.0

}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Missing parameter, provide file name!")
		return
	}
	xmlFile, err := os.Open(os.Args[1] /*"TEST.kml"*/)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer xmlFile.Close()

	xmlData, _ := ioutil.ReadAll(xmlFile)

	var k kml
	var longtitude1, latitude1, NearestLineCoNextx float64
	var d, long1, lati, NearestLineCoNextPmx, NearestLineCoInitialx, NearestLineCoInitialy float64
	WindMillDistance := []float64{}
	ActualDistance := []float64{}
	AerialDistance := []float64{}
	WindMillPoint := []string{}
	arr := []string{}
	flag := 0
	negative := 0
	i := 0
	j := 0
	dis := 0.0
	Small := 0.0

	xml.Unmarshal(xmlData, &k)

	jsonData, err := json.Marshal(k)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Write to JSON file
	jsonFile, err := os.Create("./TEST.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)

	jsonFile.Close()

	filePath := "./TEST.json"
	fmt.Printf("// reading file %s\n", filePath)
	file, err1 := ioutil.ReadFile(filePath)
	if err1 != nil {
		fmt.Printf("// error while reading file %s\n", filePath)
		fmt.Printf("File error: %v\n", err1)
		os.Exit(1)
	}

	err2 := json.Unmarshal(file, &k)
	if err2 != nil {
		fmt.Println("unmarshalling error:", err2)
		os.Exit(1)
	}
	s, er := geo.HandleWithSQL()

	if er != nil {
		panic(er)
	}

	for i = 0; i < len(k.Document.Placemark); i++ {

		arr := splitLink(k.Document.Placemark[i].LineString.Coordinates)

		if len(arr) > 0 {

			for j = 0; j < len(arr)-1; j = j + 2 {

				s.SqlConn.Exec(fmt.Sprintf("INSERT INTO points(lat, lng) VALUES(%s, %s);", arr[j], arr[j+1]))

			}
		}

	}
	//Fetching all windmill points at same time
	for i = 0; i < len(k.Document.Folder); i++ {

		for j = 0; j < len(k.Document.Folder[i].Placemark); j++ {

			if p := k.Document.Folder[i].Placemark[j].Point.Coordinates; p != "" {

				if name := k.Document.Folder[i].Placemark[j].Name; name == "WTG2100KW-S11X_PROPOSED" {

					WindMillString := k.Document.Folder[i].Placemark[j].Point.Coordinates
					Temp := splitStr(WindMillString)
					WindMillPoint = append(WindMillPoint, Temp[0])
					WindMillPoint = append(WindMillPoint, Temp[1])

				}

			}
		}
	}

	//Ascending order of windmill

	for i = 1; i < len(WindMillPoint)-1; i = i + 2 {

		for j = i + 2; j < len(WindMillPoint)-1; j = j + 2 {

			if WindMillPoint[i] > WindMillPoint[j] {

				a := WindMillPoint[i-1]
				t := WindMillPoint[i]
				WindMillPoint[i-1] = WindMillPoint[j-1]
				WindMillPoint[i] = WindMillPoint[j]
				WindMillPoint[j-1] = a
				WindMillPoint[j] = t
			}

		}

	}

	//Fetching all windmill points one by one

	for i = 0; i < 10; i = i + 2 {

		p1, err := strconv.ParseFloat(WindMillPoint[i], 64)
		p2, err := strconv.ParseFloat(WindMillPoint[i+1], 64)
		if err != nil {
			panic(err)
		}
		// fmt.Println("WindMillPoint:", p1, p2)

		WindMillDistance = WindMillDistance[:0]
		//Compare distance windmill to windmill
		for j = i; j < len(WindMillPoint)-1; j = j + 2 {

			p3, _ := strconv.ParseFloat(WindMillPoint[j], 64)
			p4, _ := strconv.ParseFloat(WindMillPoint[j+1], 64)

			if p3 != p1 && p4 != p2 {

				d = Distance(p2, p1, p4, p3)
				WindMillDistance = append(WindMillDistance, d)
				WindMillDistance = append(WindMillDistance, p3)
				WindMillDistance = append(WindMillDistance, p4)

			}

		}

		Small = WindMillDistance[0]
		long1 = WindMillDistance[1]
		lati = WindMillDistance[2]
		for k := 3; k < len(WindMillDistance); k = k + 3 {
			if Small > WindMillDistance[k] {
				Small = WindMillDistance[k]
				long1 = WindMillDistance[k+1]
				lati = WindMillDistance[k+2]
			}
		}
		// fmt.Println("Nearest windmill", Small, "meters", long1, lati)
		//initial Windmill radius

		_, NearestLineCoInitialx, NearestLineCoInitialy = NearestLinestringCo(s, p1, p2)

		//nearest windmill radius for linestring

		_, NearestLineCoNextx, _ = NearestLinestringCo(s, long1, lati)

		//Total distance between one windmill to another
		dis = 0.0
		for l := 0; l < len(k.Document.Placemark); l++ {

			arr = arr[:0]
			negative = 0
			arr = splitLink(k.Document.Placemark[l].LineString.Coordinates)

			if NearestLineCoInitialy-lati < 0 {
				negative = 1
			}

			if len(arr) > 0 {

				for j = 0; j < len(arr)-1; j = j + 2 {

					temp1, _ := strconv.ParseFloat(arr[j], 64)

					if temp1 == NearestLineCoInitialx {

						temp2, _ := strconv.ParseFloat(arr[j+1], 64)

						check, _ := strconv.ParseFloat(arr[j+3], 64)

						if j+2 < len(arr) && temp2-check < 0 && negative == 1 {

							for c := j + 2; c < len(arr)-1; c = c + 2 {

								longtitude1, _ = strconv.ParseFloat(arr[c], 64)
								latitude1, _ = strconv.ParseFloat(arr[c+1], 64)
								d = Distance(temp2, temp1, latitude1, longtitude1)
								dis = dis + d

								if longtitude1 == NearestLineCoNextx {
									flag = 0
									break
								} else {

									flag = 1
								}
								temp2 = latitude1
								temp1 = longtitude1
							}

						} else {

							for c := j - 2; c >= 0; c = c - 2 {

								longtitude1, _ = strconv.ParseFloat(arr[c], 64)
								latitude1, _ = strconv.ParseFloat(arr[c+1], 64)
								d = Distance(temp2, temp1, latitude1, longtitude1)
								dis = dis + d

								if longtitude1 == NearestLineCoNextx {
									flag = 0
									break
								} else {

									flag = 1
								}
								temp2 = latitude1
								temp1 = longtitude1
							}

						}

					}

				} //move to next placemark
				if flag == 1 {

					_, NearestLineCoNextPmx, _ = NearestLinestringCo(s, longtitude1, latitude1)

					for x := 0; x < len(k.Document.Placemark); x++ {

						arr = arr[:0]
						arr = splitLink(k.Document.Placemark[x].LineString.Coordinates)
						if len(arr) > 0 {

							for j = 0; j < len(arr)-1; j = j + 2 {

								temp1, _ := strconv.ParseFloat(arr[j], 64)

								if temp1 == NearestLineCoNextPmx {

									temp2, _ := strconv.ParseFloat(arr[j+1], 64)

									if j+2 < len(arr) {

										check, _ := strconv.ParseFloat(arr[j+3], 64)

										if temp2-check < 0 && negative == 1 {
											for c := j + 2; c < len(arr)-1; c = c + 2 {

												longtitude1, _ = strconv.ParseFloat(arr[c], 64)
												latitude1, _ = strconv.ParseFloat(arr[c+1], 64)
												d = Distance(temp2, temp1, latitude1, longtitude1)
												dis = dis + d

												if longtitude1 == NearestLineCoNextx {
													flag = 0
													break
												} else {

													flag = 1
												}
												temp2 = latitude1
												temp1 = longtitude1
											}
										} else {

											for c := j - 2; c >= 0; c = c - 2 {

												longtitude1, _ = strconv.ParseFloat(arr[c], 64)
												latitude1, _ = strconv.ParseFloat(arr[c+1], 64)
												d = Distance(temp2, temp1, latitude1, longtitude1)
												dis = dis + d

												if longtitude1 == NearestLineCoNextx {
													flag = 0
													break
												} else {

													flag = 1
												}
												temp2 = latitude1
												temp1 = longtitude1
											}
										}
									} else {

										for c := j - 2; c >= 0; c = c - 2 {

											longtitude1, _ = strconv.ParseFloat(arr[c], 64)
											latitude1, _ = strconv.ParseFloat(arr[c+1], 64)
											d = Distance(temp2, temp1, latitude1, longtitude1)
											dis = dis + d

											if longtitude1 == NearestLineCoNextx {
												flag = 0
												break
											} else {

												flag = 1
											}
											temp2 = latitude1
											temp1 = longtitude1
										}

									}
								}
							}

						}
					}

				}

			}

		}
		// fmt.Println("windmill distance from", WindMillName, "to", NextWindMillName, ":", dis, "Meters")
		ActualDistance = append(ActualDistance, dis)
		AerialDistance = append(AerialDistance, Small)

	}
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("|  Source     |   Destination  |   Aerial distance   |  Actual distance |")
	fmt.Println("-------------------------------------------------------------------------")
	for i = 0; i < len(ActualDistance); i = i + 1 {
		fmt.Println("|          ", i, "|              ", i+1, "|", AerialDistance[i], "|", ActualDistance[i], "|")

	}
	fmt.Println("-------------------------------------------------------------------------")

	FlushTestDB(s)
}
