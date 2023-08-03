package httpserver

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"regexp"
)

type PatternHandlerMap map[string]func(http.ResponseWriter, *http.Request)

//func BindMatchHandlers(patternHandlerMap PatternHandlerMap, w http.ResponseWriter, r *http.Request) {
//    for pattern, handler := range patternHandlerMap {
//        if match, err := regexp.MatchString(pattern, r.URL.Path); match && err == nil {
//            handler(w, r)
//        }
//    }
//}

type stat interface {
	Stat() (os.FileInfo, error)
}

type size interface {
	Size() int64
}

type SaveAsFunc func(file multipart.File, header *multipart.FileHeader) string

func SaveResponseAs(w http.ResponseWriter, r *http.Request, saveAsFunc SaveAsFunc) {
	if written, err := SaveAs(r, saveAsFunc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "upload file size: %d", written)
	}
}

func SaveAs(r *http.Request, saveAs SaveAsFunc) (int64, error) {
	file, header, err := r.FormFile("fieldNameHere")
	if err != nil {
		return 0, err
	}

	written, err := save(file, header, saveAs)
	if err != nil {
		return written, err
	}

	return written, nil
}

func save(file multipart.File, header *multipart.FileHeader, saveAsFunc SaveAsFunc) (int64, error) {
	dst, err := os.Create(saveAsFunc(file, header))
	if err != nil {
		return 0, err
	}

	written, err := io.Copy(dst, file)
	defer dst.Close()
	return written, nil
}

func GetSize(file multipart.File) int64 {
	if sizeInterface, ok := file.(size); ok {
		return sizeInterface.Size()
	}
	return 0
}

func ListenAndServe(addr string, patternHandlerMap PatternHandlerMap) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for pattern, handler := range patternHandlerMap {
			if match, err := regexp.MatchString(pattern, r.URL.Path); match && err == nil {
				handler(w, r)
			}
		}
	})
	if err := http.ListenAndServe(addr, mux); err != nil {
		return nil, err
	}
	return mux, nil
}
