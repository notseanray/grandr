package main
import (
    "fmt"
    "os/exec"
    "strings"
 //    "fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/widget"
)

type Screen struct {
    connected bool
    primary bool
    modes []string
    width uint32
    height uint32
}

func contains[K comparable](v K, items []K) bool {
    for _, item := range items {
        if item == v {
            return true
        }
    }
    return false
}

func fetchAndParse() string {
    cmd := exec.Command("xrandr");
    output, err := cmd.Output();
    if err != nil {
        fmt.Println("fail")
    }
    var screens = make(map[string]Screen)
    outputStr := strings.Split(string(output), "\n")
    latest := ""
    for _, i := range outputStr {
        fmt.Println(i);
        if len(i) < 16 {
            continue;
        }
        if i[0:1] == " " {
            if latest != "" {
                // screens[latest].modes
            }
            continue;
        }
        line := strings.Split(i, " ");
        if len(line) < 2 {
            continue;
        }
        fmt.Println("display: ", line[0])
        latest = line[0]
        // sep := contains[string]("primary", line);
        s := []string{}
        screens[line[0]] = Screen {
            connected: line[1] == "connected",
            primary: line[2] == "primary",
            modes: s,
        }
    }
    return "test"
}

func main() {
    fmt.Println("test");
    fetchAndParse();
 //    a := app.New()
	// w := a.NewWindow("Hello")
	//
	// hello := widget.NewLabel("Hello Fyne!")
	// w.SetContent(container.NewVBox(
	// 	hello,
	// 	widget.NewButton("Hi!", func() {
	// 		hello.SetText("Welcome :)")
	// 	}),
	// ))
	//
	// w.ShowAndRun()
}
