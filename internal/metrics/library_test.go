package metrics

import (
	"net/http"
	"testing"
	"time"

	"TOomaAh/emby_exporter_go/pkg/logger"
)

func TestLibraryCollector(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/Library/VirtualFolders/Query", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"Items":[
			{"Name":"Films","ItemId":"1","LibraryOptions":{"ContentType":"movies"}},
			{"Name":"Series","ItemId":"2","LibraryOptions":{"ContentType":"tvshows"}},
			{"Name":"Broken","ItemId":"3","LibraryOptions":{"ContentType":"movies"}},
			{"Name":"Music","ItemId":"4","LibraryOptions":{"ContentType":"music"}}
		]}`))
	})
	mux.HandleFunc("/Items", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		id := r.URL.Query().Get("ParentId")
		if id == "3" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(`{"TotalRecordCount":` + id + `0}`))
	})

	c := NewLibraryCollector(newFakeEmby(t, mux), logger.New("error"))

	start := time.Now()
	got := gather(t, c)["emby_media_item"]
	elapsed := time.Since(start)

	want := map[string]float64{"Films": 10, "Series": 20, "Music": 40}
	if len(got) != len(want) {
		t.Fatalf("got %d metrics, want %d (the failing library must be skipped): %+v", len(got), len(want), got)
	}
	for _, s := range got {
		name := s.labels["name"]
		if v, ok := want[name]; !ok || v != s.value {
			t.Errorf("emby_media_item{name=%q} = %v, want %v", name, s.value, v)
		}
	}

	// 4 requests of 200 ms: sequential would take 800 ms.
	if elapsed > 600*time.Millisecond {
		t.Errorf("collect took %s, library sizes are not requested in parallel", elapsed)
	}
}

func TestLibraryCollectorLibraryListError(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/Library/VirtualFolders/Query", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	got := gather(t, NewLibraryCollector(newFakeEmby(t, mux), logger.New("error")))
	if len(got) != 0 {
		t.Errorf("want no metric when the library list fails, got %+v", got)
	}
}
