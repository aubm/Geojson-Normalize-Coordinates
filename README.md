## Usage

Accepts a geojson file containing a single feature with a polygon geometry.
If the feature is on the date line, it will output the same feature with all the coordinates mapped to their positive equivalent.

For example, if `input.json` has content:

```json
{
   "type":"Feature",
   "properties":{

   },
   "geometry":{
      "type":"Polygon",
      "coordinates":[
         [
            [
               179,
               62
            ],
            [
               179,
               66
            ],
            [
               -170,
               66
            ],
            [
               -170,
               62
            ],
            [
               179,
               62
            ]
         ]
      ]
   }
}
```

Calling `normalize-coordinates -input-file-path=input.json` will output:

```json
{
   "type":"Feature",
   "properties":{

   },
   "geometry":{
      "type":"Polygon",
      "coordinates":[
         [
            [
               179,
               62
            ],
            [
               179,
               66
            ],
            [
               190,
               66
            ],
            [
               190,
               62
            ],
            [
               179,
               62
            ]
         ]
      ]
   }
}
```

## Installation

- From source with Go installed on the machine

```
git clone https://github.com/aubm/Geojson-Normalize-Coordinates.git $GOPATH/src/github.com/aubm/normalize-coordinates
go install github.com/aubm/normalize-coordinates
```

- Or download the appropriate release from Github from [here](https://github.com/aubm/Geojson-Normalize-Coordinates/releases)
