package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os" // Biblioteca para abrir o arquivo e trazer ele em
	"strconv"
	"strings"
	"time"
)

func main() {
	for {
		showMenu()
		action()
	}
}

func showMenu() {
	fmt.Println("1 - MONITORAR")
	fmt.Println("2 - LOGS")
	fmt.Println("3 - EXIT")
}

func action() {
	var option int
	fmt.Scan(&option)

	switch option {
	case 1:
		fmt.Println("MONITOR ON...")
		monitor()
	case 2:
		fmt.Println("[LOG] - ACCESS LOG WITH SUCESS!")
		showLogs()
	case 3:
		fmt.Println("CLOSING...")
		os.Exit(0)
	default:
		fmt.Println("COMMAND INVALID! CLOSSING...")
		os.Exit(-1)
	}
}

func monitor() {
	sites := readSitesFromArchive()
	for i := 0; i < 3; i++ {
		fmt.Println("Processing...", "[", i+1, "/", 3, "]")
		for _, site := range sites {
			testSite(site)
		}
		time.Sleep(5 * time.Second)
	}
}

func testSite(site string) {
	resp, err := http.Get(site)
	if err != nil {
		fmt.Println("[ERROR] -", err)
	}
	if resp.StatusCode == 200 {
		fmt.Println("[SUCESS] - ", resp.StatusCode, "-", site)
		registerLogs(site, true)
	} else {
		fmt.Println("[ERROR] - ", resp.StatusCode, "-", site)
		registerLogs(site, false)
	}
}

func readSitesFromArchive() []string {
	var sites []string
	archive, err := os.Open("sites.txt") // Devolve o espaço na memoria que se encontra
	// archive, err := os.ReadFile("sites.txt") // Esse devolve um array de bytes e readFile normalmente utilizado para ler o arquivo inteiro de uma vez

	// Quando queremos tratar erros, basta ver se são diferentes de <nil> que seria o null em outras linguagens
	if err != nil {
		fmt.Println("[ERROR] -", err)
		os.Exit(-1)
	}
	// fmt.Println(string(archive)) //Casting explicito de array de bytes para string

	reader := bufio.NewReader(archive) //bufio nos entrega uma classe que consegue ler um arquivo pelo seu endereço de memoria que seria o Reader
	for {
		line, err := reader.ReadString('\n') // Usamos aspas simples para passar o byte limitador para ler o arquivo.
		line = strings.TrimSpace(line)
		sites = append(sites, line)
		if err == io.EOF {
			break // Usamos o break para parar o for ja que esse erro de io.EOF é com relação que não existe mais linhas para se ler
		}
	}
	archive.Close() // sempre é bom fechar o arquivo depois manusear ele.
	fmt.Println(sites)
	return sites
}

func registerLogs(site string, status bool) {
	/*
		Essa funcao OpenFIle é mais poderosa que a Open que usamos
		para ler, porque com ela, criamos arquivos e tambem podemos passar
		flag se queremos que arquviso seja apenas de leitura ou escrita
		etc...
	*/
	archive, err := os.OpenFile("log.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}
	/*
		Ler a doc para entender o format do time.now
		https://go.dev/src/time/format.go
	*/
	archive.WriteString("[" + time.Now().Format("02/01/2006 15:04:05") + "]" +
		"[INFO]" + site + " - online: " + strconv.FormatBool(status) + "\n")
	archive.Close()
}

func showLogs() {
	archive, err := os.ReadFile("log.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(archive))
}
