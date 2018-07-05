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
-------------------------------------------------------------------------
|  Source     |   Destination  |   Aerial distance   |  Actual distance |
-------------------------------------------------------------------------
|           0 |               1 | 410.68875773780644 | 458.0948004254824 |
|           1 |               2 | 612.6801424949972 | 640.5472693065335 |
|           2 |               3 | 929.4247017776966 | 1279.4511562583766 |
|           3 |               4 | 465.41452138347773 | 495.38087872066933 |
|           4 |               5 | 702.6446780077495 | 824.9538595726046 |
-------------------------------------------------------------------------
```


