// How do you convert a string to a byte array?
[]byte("string")

// How do you print Unix epoch number as string?
import "fmt"
fmt.Sprintf("%d", time.Now().Unix())

// How do you write a string to a file?
import "ioutil"
ioutil.WriteFile(path, []byte("string"), 0644)

// How do you create a directory?
import "os"
os.Mkdir(folder, 0755)

// How do you read a whole file?
import "ioutil"
content, error := ioutil.ReadFile(path)

// How do you sleep for x seconds?
import "time"
time.Sleep(time.Duration(3) * time.Second)

// How do you parse a JSON string with an array of ints?
import "encoding/json"
var ids []int
err := json.Unmarshal([]byte(jsonString), &ids)

// Can you convert to JSON a map[int]string ?
// No, keys in maps can only be marshalled/un-marshalled if they are strings. map[string]string works.

// How do you check if a map has a certain key?
if _, ok := map[key]; ok {}

// How do you sort an array?
// Define a custom type:
type ByTime []os.FileInfo
// Implement 3 methods:
func (a ByTime) Len() int { return len(a) }
func (a ByTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].ModTime().Unix() > a[j].ModTime().Unix() }
// Finally use like this:
sort.Sort(ByTime(fileList))

// How do you capture Control C before exiting?
import "os"
c := make(chan os.Signal, 1)
signal.Notify(c, os.Interrupt)
go func() {
    for range c {
        //do something here
        os.Exit(0)
}()

// How do you delete a file?
import "os"
os.Remove(path) error

// How do you define command line flags?
import "flag"
var (
    port = flag.String("port", "8080", "web server port")
    static = flag.String("static", "./static/", "static folder")
    config = flag.String("config", "config.json", "path to config with all credentials")
)
flag.Parse()
