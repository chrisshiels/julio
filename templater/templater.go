// 'templater.go'.
// Chris Shiels.


package main


import (
    "fmt"
    "io/ioutil"
    "os"
    "strings"
    "text/template"
)


const exitsuccess = 0
const exitfailure = 1


func dotstodashes(s string) string {
    return strings.Replace(s, ".", "-", -1)
}


func templater(filename string, vars map[string]string) error {
    file, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer file.Close()

    templatebytes, err := ioutil.ReadAll(file)
    if err != nil {
        return err
    }
    templatestring := string(templatebytes)

    t, err := template.New(filename).
                       Funcs(template.FuncMap { "dotstodashes": dotstodashes }).
                       Parse(templatestring)
    if err != nil {
        return err
    }

    err = t.Execute(os.Stdout, vars)
    if err != nil {
        return err
    }

    return nil
}


func main() {
    vars := make(map[string]string)

    if len(os.Args[1:]) == 0 {
        fmt.Fprintln(os.Stderr,
                     "Usage:  templater [ key=value ... ] filename.template")
        os.Exit(exitsuccess)
    }

    for _, arg := range os.Args[1:len(os.Args) - 1] {
        keyvalue := strings.SplitN(arg, "=", 2)
        if len(keyvalue) == 2 {
            vars[keyvalue[0]] = keyvalue[1]
        } else {
            vars[keyvalue[0]] = "1"
        }
    }

    filename := os.Args[len(os.Args) - 1]

    if err := templater(filename, vars); err != nil {
        fmt.Fprintf(os.Stderr, "templater: %s\n", err)
        os.Exit(exitfailure)
    }

    os.Exit(exitsuccess)
}
