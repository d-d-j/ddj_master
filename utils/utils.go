package utils

import (
	"log"
)

func Log(v ...interface{}) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.Println(v...)
}
