# kmllinetracing

Currently,  kmllinetracing provides the following functionality:

-kmltojson parsing
-finding nearest windmill point 
-finding nearest line string coordinate using circle radius formula with sql
-4 windmill's placemark traced

## Running the code

To run the code from Linetracing directory simply execute:

```
DB=mysql GO_ENV=test go run LineTracing.go TEST.kml
```

And the output will be something like this:

```
// reading file ./TEST.json
WindMillPoint: 85.2934715082 25.0739799545
Nearest windmill 410.68875773780644 meters 85.2937873022 25.0776581512
Nearest windmill distance is  : 458.0948004254824 Meters
WindMillPoint: 85.2937873022 25.0776581512
Nearest windmill 612.6801424949972 meters 85.2943755681 25.0831361328
Nearest windmill distance is  : 640.5472693065335 Meters
WindMillPoint: 85.3277935074 25.0781265342
Nearest windmill 929.4247017776966 meters 85.3331397194 25.0713249995
Nearest windmill distance is  : 1279.4511562583766 Meters
WindMillPoint: 85.2943755681 25.0831361328
Nearest windmill 465.41452138347773 meters 85.2950022124 25.0872783475
Nearest windmill distance is  : 495.38087872066933 Meters
WindMillPoint: 85.2950022124 25.0872783475
Nearest windmill 702.6446780077495 meters 85.2950733563 25.0935900198
Nearest windmill distance is  : 824.9538595726046 Meters
```


