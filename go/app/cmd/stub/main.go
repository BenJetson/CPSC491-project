package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func makeHandle(tier string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		out := struct {
			Error string `json:"error"`
		}{
			fmt.Sprintf(
				"You have successfully hit the %s tier of the API server. "+
					"Unfortunately, this application is a stub and "+
					"does nothing.",
				tier,
			),
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")

		e := json.NewEncoder(w)
		err := e.Encode(out)
		if err != nil {
			log.Panicf("could not encode response, received: %v\n", err)
		}
	}
}

func main() {
	tier := os.Getenv("TIER")
	if len(tier) < 1 {
		log.Fatalln("TIER cannot be blank")
	}

	portStr := os.Getenv("PORT")
	if len(portStr) < 1 {
		log.Fatalln("PORT cannot be blank")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("PORT not an integer, received: %v\n", err)
	} else if port < 1000 {
		log.Fatalf("PORT out of range, must be > 1000")
	}

	h := makeHandle(tier)
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = s.ListenAndServe()
	log.Fatalln(err)
}
