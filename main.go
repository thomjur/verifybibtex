package main

import (
	"fmt"

	"github.com/thomjur/verifybibtex/parser"
)

func main() {
	fmt.Println("Juten Tach!")
	//entry := `@article{id1234,author={Jurczyk, Thomas},date={20.12.2023}}`
	entry := `@article{schmidt2024,author = {Schmidt, Anna and Müller, Bernd and {O'Connor}, Claire and García, Diego},language = "Deutsch"}`

	// 	entry := `
	// 	@book{schmidt2024,
	//   author       = {Schmidt, Anna and Müller, Bernd and {O'Connor}, Claire and García, Diego},
	//   editor       = "Weber, Eva and {D'Amico}, Fabio",
	//   title        = {Fortgeschrittene Datenanalyse (mit, = in dem Text) mit Python: Methoden und Anwendungen}, publisher    = {Technik Verlag}  ,
	//   year         = {2024},
	//   volume       = {3},
	//   series       = {Datenwissenschaftliche Studien},
	//   address      = {München},
	//   edition      = {2., überarbeitete und erweiterte Auflage},
	//   month        = {März},
	//   isbn         = {978-3-16-148410-0},
	//   doi          = {10.1000/182},
	//   url          = {https://www.technik-verlag.de/buecher/fortgeschrittene-datenanalyse},
	//   note         = {Beinhaltet ein Kapitel über maschinelles Lernen},
	//   abstract     = {Dieses Buch bietet eine umfassende Einführung in fortgeschrittene Methoden der Datenanalyse mit Python, einschließlich praxisnaher Anwendungen und Fallstudien.},
	//   keywords     = {Datenanalyse, Python, maschinelles Lernen, Statistik},
	//   language     = {Deutsch}
	// }

	// 	`

	parsedEntry, err := parser.ParseNewEntry(entry)

	if err != nil {
		fmt.Println(err.Error())
	}
	for k, v := range parsedEntry.Fields {
		fmt.Println(k, "::", v)
	}

}
