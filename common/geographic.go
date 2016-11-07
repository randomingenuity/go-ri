package ricommon

import (
    "fmt"

    "google.golang.org/appengine"

    "github.com/gansidui/geohash"
    "github.com/dsoprea/go-logging"
)

// Constants
const (
    MaxGeohashPrecision = 12
)

func GetBoundingHashPrefixForBox(ob *geohash.Box, center *appengine.GeoPoint) (hash string, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    for precision := MaxGeohashPrecision; precision > 0 ; precision-- {
        h, ib := geohash.Encode(center.Lat, center.Lng, precision)
        
        // Loop as long as any part of the inner geohash box remains inside the 
        // outer bounding box (we want to grow the inner box until it fully 
        // exceeds the outer box).
        if ib.MinLat < ob.MinLat || ib.MaxLat < ob.MaxLat || ib.MinLng < ob.MinLng || ib.MaxLng < ob.MaxLng {
            continue
        }

        return h, nil
    }

    // The bounding coordinates were probably nonsense.
    err = fmt.Errorf("box is invalid")
    log.Panic(err)
    return "", err
}

func EncodeCoordinatesToGeohash(latitude, longitude float64) (hash string, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = state.(error)
        }
    }()

    hash, _ = geohash.Encode(latitude, longitude, MaxGeohashPrecision)
    isValid := false

    for _, char := range hash {
        if char != 'z' {
            isValid = true
            break
        }
    }

    if isValid == false {
        log.Panic(fmt.Errorf("could not encode coordinates: (%f, %f)", latitude, longitude))
    }

    return hash, nil
}