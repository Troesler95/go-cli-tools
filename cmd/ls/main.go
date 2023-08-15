package main

import (
	"flag"
	"fmt"
	"io/fs"
	"math"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"

	. "github.com/troesler95/go-cli-tools"
	"golang.org/x/term"
)

const (
	MaxTerminalColumns int = 12
)

type Options struct {
	UseColor bool
	PrintAll bool
	ListLong bool
}

func setupCliFlags() (opts Options, args []string) {
	flag.Usage = func() {
		write := flag.CommandLine.Output()
		fmt.Fprintln(write, "usage: ls [OPTIONS]... [DIRECTORY]")
		fmt.Fprintln(write, "List information about the [DIRECTORY]. If no directory is specified, the current directory is used.")

		fmt.Fprintln(write, "")
		flag.PrintDefaults()
	}

	printAll := false
	useColor := flag.Bool("color", true, "toggles console output colors")
	listLong := flag.Bool("l", false, "use a long listing format")
	flag.BoolVar(&printAll, "all", false, "do not ignore entries staritng with .")
	flag.BoolVar(&printAll, "a", false, "do not ignore entries staritng with .")

	flag.Parse()

	if flag.NArg() > 1 {
		panic(fmt.Sprintln("expected at most 1 argument but got ", flag.NArg()))
	}

	return Options{
		UseColor: *useColor,
		PrintAll: printAll,
		ListLong: *listLong,
	}, flag.Args()
}

var fileInfoTypeColorMap map[string]Color = map[string]Color{
	"": Color{},
}

func main() {
	opts, args := setupCliFlags()

	var printer *ColorizedPrinter
	if opts.UseColor {
		printer = NewColorizedPrinter(NewColor(FgDefault, BgDefault, DefaultText))
	} else {
		printer = NewDefaultPrinter()
	}

	// TODO: check if stdout is terminal, othewise default to some value
	terminalWidth := 200
	var err error
	if term.IsTerminal(int(os.Stdin.Fd())) {
		terminalWidth, _, err = term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
	}

	var wd string
	if len(args) == 0 {
		wd, err = os.Getwd()
		if err != nil || len(wd) == 0 {
			printer.PrintError(err, "unable to determine pwd")
			panic(1)
		}
	} else {
		wd = args[0]
	}

	dirEntries, err := os.ReadDir(wd)
	if err != nil {
		printer.PrintError(err, "unable to read directory")
		panic(1)
	}

	if opts.ListLong {
		listLong(opts, printer, dirEntries)
		os.Exit(0)
	}

	longestStrLen := 0
	for _, entry := range dirEntries {
		if len(entry.Name()) > longestStrLen {
			longestStrLen = len(entry.Name())
		}
	}

	maxColumns := int(math.Min(float64(MaxTerminalColumns), math.Floor(float64(terminalWidth)/float64(longestStrLen+2))))
	colNum := 0
	for idx, entry := range dirEntries {
		if !opts.PrintAll && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		if opts.UseColor && entry.IsDir() {
			printer.PrintColor(NewColor(FgBlue, BgDefault, DefaultText), entry.Name(), "  ")
		} else {
			printer.Print(entry.Name(), "  ")
		}
		colNum++

		if idx+1 == maxColumns {
			printer.Println()
		}
	}

	printer.Println()
}

func listLong(opts Options, printer *ColorizedPrinter, dirEntries []fs.DirEntry) {
	blockCount := int64(0)
	output := strings.Builder{}
	for _, obj := range dirEntries {
		if !opts.PrintAll && strings.HasPrefix(obj.Name(), ".") {
			continue
		}

		info, err := obj.Info()
		if err != nil {
			printer.PrintError(err, fmt.Sprintf("unable to get FileInfo for %s", obj.Name()))
			panic(1)
		}

		modTime := info.ModTime()
		yearThreshold := time.Now().AddDate(-1, 0, 0)
		modTimeStr := modTime.Format("Jan 02 15:04")
		if modTime.Before(yearThreshold) {
			modTimeStr = modTime.Format("Jan 02  2006")
		}

		ownerString := "?"
		groupString := "?"
		numHardLinks := 1
		// everything below is linux specific
		if sys := info.Sys(); sys != nil {
			if stat, ok := sys.(*syscall.Stat_t); ok {
				numHardLinks = int(stat.Nlink)

				u, err := user.LookupId(strconv.Itoa(int(stat.Uid)))
				if err != nil {
					ownerString = strconv.Itoa(int(stat.Uid))
				} else {
					ownerString = u.Username
				}

				// below calculation is adapted from posix calculations defined in
				// https://github.com/rofl0r/gnulib/blob/15af560b45353d42e7553f3b98469290001b5ff6/lib/stat-size.h#L60
				//
				// the * 4 at the end is a hack, unsure if that will work in all scenarios
				blockCount += (stat.Size/stat.Blksize + int64(btoi(stat.Size%stat.Blksize != 0))) * 4

				g, err := user.LookupGroup(strconv.Itoa(int(stat.Gid)))
				if err == nil {
					groupString = g.Name
				} else if u != nil {
					groupString = u.Username
				} else {
					groupString = strconv.Itoa(int(stat.Gid))
				}
			}
		}

		fileName := info.Name()
		if info.IsDir() {
			fileName = Colorize(fileName, Color{Foreground: FgBlack, Background: BgDefault})
		}

		color := NewDefaultColor()
		if opts.UseColor && info.IsDir() {
			color = NewColor(FgBlue, BgDefault, DefaultText)
		}

		printer.Fprintf(&output, "%s %d %s %s %4d %13s  ",
			info.Mode(),
			numHardLinks,
			ownerString,
			groupString,
			info.Size(),
			modTimeStr,
		)
		printer.FprintlnColor(&output, color, fileName)
	}
	fmt.Println("total ", blockCount)
	fmt.Println(output.String())
}

func btoi(boolVal bool) int {
	if boolVal {
		return 1
	}
	return 0
}
