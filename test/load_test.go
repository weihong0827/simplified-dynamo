package main

import (
    "net/http"
    "sync"
    "testing"
    "time"
    "fmt"
)

func TestLoadPut(t *testing.T) {
    start := time.Now()

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            puturl := "http://127.0.0.1:8080/put?key=foo" + fmt.Sprint(i) + "&value=bar"
            putresp, err := http.NewRequest(http.MethodPut, puturl, nil)
            if err != nil {
                t.Fatal(err)
            }

            client := &http.Client{}
            resp, err := client.Do(putresp)
            if err != nil {
                t.Fatal(err)
            }

            if resp.StatusCode != http.StatusOK {
                t.Errorf("Expected status OK, got %v", resp.Status)
            }
        }(i)
    }
    wg.Wait()

    t.Logf("TestLoadPut 10 Put took %v", time.Since(start))
}

// load test for get
func TestLoadGet(t *testing.T) {
    start := time.Now()

    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            geturl := "http://127.0.0.1:8080/get?key=foo" + fmt.Sprint(i)
            resp, err := http.Get(geturl)
            if err != nil {
                t.Fatal(err)
            }

            if resp.StatusCode != http.StatusOK {
                t.Errorf("Expected status OK, got %v", resp.Status)
            }
        }(i)
    }
    wg.Wait()

    t.Logf("TestLoadGet 10 Get took %v", time.Since(start))
}