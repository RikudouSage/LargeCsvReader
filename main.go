package main

import (
	"LargeCsvReader/browser"
	"LargeCsvReader/widgets"
	"embed"
	"encoding/csv"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/hashicorp/go-version"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
)

const repository = "https://github.com/RikudouSage/LargeCsvReader"

//go:embed translation
var translations embed.FS

//go:embed assets/appversion
var appVersion string

func checkForNewVersion(window fyne.Window) {
	currentVersion, err := version.NewVersion(strings.TrimSpace(appVersion))
	if err != nil || currentVersion.String() == "dev" {
		fmt.Println(err)
		return
	}

	const url = repository + "/releases/latest"
	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	response, err := httpClient.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	response.Body.Close()

	redirectUrl := response.Header.Get("Location")
	if redirectUrl == "" {
		return
	}

	regex := regexp.MustCompile("https://.*/v([0-9.]+)")
	matches := regex.FindStringSubmatch(redirectUrl)
	if len(matches) < 2 {
		return
	}

	newestVersion, err := version.NewVersion(matches[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	if newestVersion.GreaterThan(currentVersion) {
		dialog.ShowConfirm(
			lang.X("app.new_version.title", "New version found!"),
			lang.X("app.new_version.description", "Do you want to download the newest version?"),
			func(result bool) {
				if !result {
					return
				}

				err = browser.OpenUrl(redirectUrl)
				if err != nil {
					fmt.Println(err)
					return
				}
			},
			window,
		)
	}
}

func showPreviewWindow(filePath string, fyneApp fyne.App) {
	window := fyneApp.NewWindow(lang.X("app.preview", "Preview"))
	window.Resize(fyne.NewSize(640, 480))

	typeToCharMap := map[string]rune{
		lang.X("app.separator.comma", "Comma"):         ',',
		lang.X("app.separator.semicolon", "Semicolon"): ';',
		lang.X("app.separator.tab", "Tab"):             '\t',
	}
	separatorOptions := make([]string, len(typeToCharMap))
	i := 0
	for key := range typeToCharMap {
		separatorOptions[i] = key
		i++
	}
	sort.Slice(separatorOptions, func(i, j int) bool {
		return separatorOptions[i] < separatorOptions[j]
	})

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	var csvFile *csv.Reader
	rows := make([][]string, 5)

	previewText := widget.NewLabel(lang.X("app.preview", "Preview"))
	previewText.Hide()

	tableContainer := container.NewStack()
	lineSeparator1 := widget.NewSeparator()
	lineSeparator2 := widget.NewSeparator()

	lineSeparator1.Hide()
	lineSeparator2.Hide()

	createTable := func() *widget.Table {
		table := widget.NewTable(
			func() (int, int) {
				return len(rows), len(rows[0])
			},
			func() fyne.CanvasObject {
				return widgets.NewCellWidget("default (hopefully) large enough text", nil) // placeholder to specify width
			},
			func(id widget.TableCellID, object fyne.CanvasObject) {
				cell := object.(*widgets.CellWidget)
				cell.SetText(rows[id.Row][id.Col])
				cell.OnRightClick = func(event *fyne.PointEvent) {
					items := []*fyne.MenuItem{
						fyne.NewMenuItem(lang.X("app.copy_to_clipboard", "Copy"), func() {
							window.Clipboard().SetContent(cell.Text)
						}),
					}
					menu := fyne.NewMenu("", items...)
					canvas := fyne.CurrentApp().Driver().CanvasForObject(cell)
					widget.ShowPopUpMenuAtPosition(menu, canvas, event.AbsolutePosition)
				}
			},
		)

		return table
	}

	var loadMoreButton *widget.Button
	loadMoreButton = widget.NewButton(lang.X("app.load_more_button", "Load more"), func() {
		for range 5 {
			line, err := csvFile.Read()
			if err == io.EOF {
				loadMoreButton.Disable()
				break
			}
			if err != nil {
				panic(err)
			}
			rows = append(rows, line)
		}

		tableContainer.Objects = []fyne.CanvasObject{createTable()}
		tableContainer.Refresh()
	})
	loadMoreButton.Hide()

	separatorSelect := widget.NewSelect(separatorOptions, func(selected string) {
		csvFile = csv.NewReader(file)
		csvFile.LazyQuotes = true
		csvFile.Comma = typeToCharMap[selected]

		_, err = file.Seek(0, 0)
		if err != nil {
			panic(err)
		}
		if err != nil {
			panic(err)
		}

		rows = [][]string{}
		for range 5 {
			line, err := csvFile.Read()
			if err == io.EOF {
				loadMoreButton.Disable()
				break
			}
			if err != nil {
				panic(err)
			}
			rows = append(rows, line)
		}

		previewText.Show()
		lineSeparator2.Show()
		lineSeparator1.Show()
		loadMoreButton.Show()

		table := createTable()
		tableContainer.Objects = []fyne.CanvasObject{table}
		tableContainer.Refresh()
	})

	mainContainer := container.NewBorder(
		container.NewVBox(
			widget.NewLabel(lang.X("app.separator", "Separator")),
			separatorSelect,
			lineSeparator1,
			container.NewHBox(
				previewText,
				loadMoreButton,
			),
			lineSeparator2,
		),
		nil, nil, nil,
		tableContainer,
	)

	window.SetContent(
		container.New(
			layout.NewCustomPaddedLayout(8, 32, 32, 32),
			mainContainer,
		),
	)
	window.SetOnClosed(func() {
		file.Close()
	})

	window.Show()
}

func main() {
	fyneApp := app.NewWithID("1af17320-d9ef-4d2b-aae1-8226b14a177a")
	err := lang.AddTranslationsFS(translations, "translation")
	if err != nil {
		panic(err)
	}

	window := fyneApp.NewWindow(lang.X("app.title", "Large CSV Reader") + " (" + strings.TrimSpace(appVersion) + ")")
	window.Resize(fyne.NewSize(640, 480))
	window.SetMaster()

	openDialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		if file != nil {
			showPreviewWindow(file.URI().Path(), fyneApp)
		}
	}, window)
	openDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))

	repoUrl, err := url.Parse(repository)
	if err != nil {
		panic(err)
	}

	window.SetContent(container.NewVBox(
		container.NewCenter(
			widget.NewLabel(lang.X("app.choose_csv_button.description", "Using the button below choose a CSV you want to open")),
		),
		container.New(
			layout.NewCustomPaddedLayout(8, 0, 32, 32),
			widget.NewButton(lang.X("app.choose_csv_button", "Open CSV"), func() {
				openDialog.Show()
			}),
		),
		layout.NewSpacer(),
		container.NewCenter(
			widget.NewHyperlink(repository, repoUrl),
		),
	))

	go checkForNewVersion(window)

	window.Show()
	fyneApp.Run()
}
