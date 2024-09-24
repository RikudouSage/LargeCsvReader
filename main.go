package main

import (
	"LargeCsvReader/widgets"
	"embed"
	"encoding/csv"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"io"
	"os"
)

//go:embed translation
var translations embed.FS

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

	window := fyneApp.NewWindow(lang.X("app.title", "Large CSV Reader"))
	window.Resize(fyne.NewSize(640, 480))
	window.SetMaster()

	openDialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		if file != nil {
			showPreviewWindow(file.URI().Path(), fyneApp)
		}
	}, window)
	openDialog.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))

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
	))

	window.Show()
	fyneApp.Run()
}
