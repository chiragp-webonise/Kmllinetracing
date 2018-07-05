# kmllinetracing

Currently,  kmllinetracing provides the following functionality:

-kmltojson parsing
-finding nearest windmill point 
-finding nearest line string coordinate using circle radius formula with sql
-4 windmill's placemark traced

## Getting Started

#### Prerequisites

First you need to clone the repo of golang geo to get the geo library

```
git clone git@github.com:chiragp-webonise/golang-geo.git
```

## Running the code

To run the code from Linetracing directory simply execute:

```
DB=mysql GO_ENV=test go run LineTracing.go TEST.kml
```

And the output will be something like this:

```
// reading file ./TEST.json
windmill distance from 1 to 2 : 458.0948004254824 Meters
windmill distance from 2 to 3 : 640.5472693065335 Meters
windmill distance from 3 to 4 : 1279.4511562583766 Meters
windmill distance from 4 to 5 : 495.38087872066933 Meters
windmill distance from 5 to 6 : 824.9538595726046 Meters
```


