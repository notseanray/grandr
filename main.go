package main

import (
    "fmt"
    "os/exec"
    "reflect"
    "sort"
    "strconv"
    "strings"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

type Connection struct {
    name      string
    connected bool
    primary   bool
    modes     []string
    width     uint32
    height    uint32
}

type RootScreen struct {
    screenNum uint8
    minimumX  int16
    minimumY  int16
    currentX  int16
    currentY  int16
    maxX      int16
    minY      int16
}

type XrandrOutput struct {
    root    RootScreen
    screens map[string]Connection
}

func contains[K comparable](v K, items []K) bool {
    for _, item := range items {
        if item == v {
            return true
        }
    }
    return false
}

type NumericInt interface {
    int | int8 | int16 | int32 | int64 |
        uint | uint8 | uint16 | uint32 | uint64
}

// this is janky, does it even work??
func convertFallible[V NumericInt](value string, defaltVal V) V {
    result, err := strconv.ParseInt(value, 10, reflect.TypeOf((*V)(nil)).Elem().Bits())
    if err != nil {
        return defaltVal
    }
    return V(result)
}

func parseMonitors() string {
    cmd := exec.Command("xrandr --listmonitors")
    output, err := cmd.Output()
    if err != nil {
        fmt.Println("fail")
    }
    outputStr := strings.Split(string(output), "\n")
    first := true
    screenCount := uint8(0)
    for _, i := range outputStr {
        line := strings.Split(i, " ")
        if first {
            first = false
            screenCount = convertFallible[uint8](line[1], 0)
        }
        fmt.Println(screenCount)
    }
    return "test"
}

func fetchAndParse() XrandrOutput {
    cmd := exec.Command("xrandr")
    output, err := cmd.Output()
    if err != nil {
        fmt.Println("fail")
    }
    var screens = make(map[string]Connection)
    outputStr := strings.Split(string(output), "\n")
    latest := ""
    first := true
    rootscreen := RootScreen{
        screenNum: 0,
        minimumX:  0,
        minimumY:  0,
        currentX:  0,
        currentY:  0,
        maxX:      0,
        minY:      0,
    }
    _ = rootscreen
    for _, i := range outputStr {
        if len(i) < 16 {
            continue
        }
        if first {
            first = false
            line := strings.Split(i, " ")
            // the format is different
            if len(line) != 14 {
                continue
            }
            rootscreen = RootScreen{
                screenNum: convertFallible[uint8](line[1][0:len(line[1])-1], 0),
                minimumX:  convertFallible[int16](line[3], 0),
                minimumY:  convertFallible[int16](line[5][0:len(line[5])-1], 0),
                currentX:  convertFallible[int16](line[7], 0),
                currentY:  convertFallible[int16](line[9][0:len(line[9])-1], 0),
                maxX:      convertFallible[int16](line[11], 0),
                minY:      convertFallible[int16](line[13][0:len(line[13])-1], 0),
            }
            fmt.Println(rootscreen)
            continue
        }
        if i[0:1] == " " {
            if latest != "" {
                // screens[latest].modes
            }
            continue
        }
        line := strings.Split(i, " ")
        if len(line) < 2 {
            continue
        }
        fmt.Println("connector: ", line[0])
        latest = line[0]
        // sep := contains[string]("primary", line);
        s := []string{}
        screens[line[0]] = Connection{
            name:      line[0],
            connected: line[1] == "connected",
            primary:   line[2] == "primary",
            modes:     s,
        }
    }
    return XrandrOutput{root: rootscreen, screens: screens}
}

func test(val string) {
    fmt.Println(val)
}

type ConnectionSorted []Connection

func score(c Connection) int {
    score := 0
    if c.primary {
        score -= 100
    }
    if c.connected {
        score -= 50
    }
    return score
}

func (s ConnectionSorted) Len() int {
    return len(s)
}

func (s ConnectionSorted) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

func (s ConnectionSorted) Less(i, j int) bool {
    return score(s[i]) < score(s[j])
}

func main() {
    data := fetchAndParse()
    sortedList := []Connection{}
    for _, d := range data.screens {
        sortedList = append(sortedList, d)
    }
    sort.Sort(ConnectionSorted(sortedList))
    fmt.Println(sortedList)
    a := app.New()
    w := a.NewWindow("grandr")
    var buttons []fyne.CanvasObject
    buttons = append(buttons, widget.NewLabel("Screen: "+strconv.Itoa(int(data.root.screenNum))))
    for _, d := range sortedList {
        // change this to have default values then reassign
        if d.connected {
            screenEntry := widget.NewButton("set as primary", func() {
                // parse xrandr --listmonitors, based on that set as primary
            })
            screenCard := widget.NewCard(
                d.name,
                "Connected: "+strconv.FormatBool(d.connected)+" Primary: "+strconv.FormatBool(d.primary), screenEntry)
            buttons = append(buttons, screenCard)
        } else {
            screenEntry := widget.NewButton("test", func() {
                // parse xrandr --listmonitors, based on that set as primary
            })
            screenCard := widget.NewCard(
                d.name,
                "Connected: "+strconv.FormatBool(d.connected), screenEntry)
            buttons = append(buttons, screenCard)
        }
    }
    content := container.NewGridWithRows(len(data.screens)+1, buttons...)
    w.SetContent(content)
    w.ShowAndRun()
}
