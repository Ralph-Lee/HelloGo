package main
 
import ("os"
	"log"
	"fmt"
	"time"
	"strconv"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"crypto/rand"
)
 
func main() {
    download_link := "http://www.xieyuluo.com/res/libo.p12"
    path, err := download_file(download_link)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(path)
}
 
func download_file(download_link string) (string, error) {
    ext := filepath.Ext(download_link)
    dir_path, err := dir_full_path()
    if err != nil {
        return "", err
    }
    file_name := rand_str(10) + ext
    file_path := dir_path + file_name
    
    os.Mkdir(dir_path, 0666)
    
    file, err := os.Create(file_path)
    if err != nil {
        return "", err
    }
    defer file.Close()
 
    res, err := http.Get(download_link)
    if err != nil {
        return "", err
    }
    defer res.Body.Close()
 
    file_content, err := ioutil.ReadAll(res.Body)
 
    if err != nil {
        return "", err
    }
 
    // returns file size and err
    _, err = file.Write(file_content)
 
    if err != nil {
        return "", err
    }
 
    return file_path, nil
}
 
func dir_full_path() (string, error) {
    path, err := filepath.Abs("files")
 
    if err != nil {
        return "", err
    }
 
    t := time.Now()
 
    s := path +
        string(os.PathSeparator) +
        strconv.Itoa(t.Day()) +
        "_" +
        strconv.Itoa(int(t.Month())) +
        "_" +
        strconv.Itoa(t.Year()) +
        string(os.PathSeparator)
 
    return s, nil
}
 
func rand_str(n int) string {
    alphanum := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
    var bytes = make([]byte, n)
    rand.Read(bytes)
    for i, b := range bytes {
        bytes[i] = alphanum[b%byte(len(alphanum))]
    }
    return string(bytes)
}
