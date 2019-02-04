package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const monitoramento = 3
const delay = 5

func main() {

	exibeIntro()

	//instrucao for para rodar indefinidamente ate usuario encerrar monitoramento com for sem regras
	for {

		exibeMenu()

		comando := leComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			imprimeLogs()
		case 3:
			fmt.Println("Saindo...")
			os.Exit(0) //sai do programinha
		default:
			fmt.Println("Esta opção não existe!")
			os.Exit(-1) //-1 para escolhas que nao existem no menu e gerar msg erro
		}

		leComando()
	}

}

//funcao de exibicao de informacao
func exibeIntro() {
	nome := "Monitoserv"
	versao := 1.1
	fmt.Println("Monitoramento de Serviços", nome)
	fmt.Println("Este sistema está na versão", versao)
}

//funcao para exibir o menu
func exibeMenu() {
	fmt.Println("1 - Iniciar Monitoramento")
	fmt.Println("2 - Exibir Logs")
	fmt.Println("3 - Saír")
}

//funcao para ler opcao menu escolhida
func leComando() int {
	var comandoLido int
	fmt.Scan(&comandoLido)
	fmt.Println("O comando escolhido foi", comandoLido, "na posicao de memória", &comandoLido)
	fmt.Println("")

	return comandoLido
}

func iniciarMonitoramento() {
	fmt.Println("Monitoramento iniciado... por favor, aguarde retorno...")

	// servicos := []string{
	// 	"http://www.4effect.com.br/",
	// 	"http://www.thiagolucio.com.br/",
	// 	"http://www.radarturismojf.com.br/",
	// 	"http://www.adelmoral.com.br/",
	// }

	//refatorando para ler de uma lista de arquivo de texto externo
	servicos := leservicosDeArquivo()

	for i := 0; i < monitoramento; i++ {
		//usando o range do slices
		for i, servico := range servicos {
			fmt.Println("Testando URI", i, " => ", servico)
			testaservicos(servico)
		}
		time.Sleep(delay * time.Second)
		fmt.Println("")
	}
}

func testaservicos(servico string) {
	resp, err := http.Get(servico) //como o Get servico retorna mais de um valor, somente pegamos a resposta(resp) e ignoramos o segundo retorno.

	if err != nil {
		fmt.Println("Ocorreu um erro:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Println("servico:", servico, "foi carregado com sucesso! - Status Code:", resp.StatusCode)
		registraLog(servico, true) //se o servico estiver no ar ele passara o servico e o status true indicando que estava no ar para o arquivo de log
	} else {
		fmt.Println("servico:", servico, "esta com problemas. Status Code:", resp.StatusCode)
		registraLog(servico, false) //se o servico estiver no ar ele passara o servico e o status false indicando que estava FORA do ar para o arquivo de log

	}
}

func leservicosDeArquivo() []string {

	var servicos []string

	//arquivo, err := ioutil.ReadFile("servicos.txt") //usa o ioutil.ReadFile para ler o arquivo (o que tem dentro dele)
	arquivo, err := os.Open("servicos.txt") // Usa o os.Open para abrir o arquivo de texto

	if err != nil {
		fmt.Println("Ocorreu o erro", err)
	}

	leitor := bufio.NewReader(arquivo)

	for {
		linha, err := leitor.ReadString('\n') //aqui ele le somente até a primeira quebra de linha
		linha = strings.TrimSpace(linha)      //usa o cortar Trim para tirar todo espaço em branco da linha. Não esta criando uma nova var apenas pegando o mesmo valor dela sem os espaços e outros \n da linha

		servicos = append(servicos, linha)

		if err == io.EOF {
			break
		}
	}

	arquivo.Close()

	return servicos
}

func registraLog(servico string, status bool) {
	//explicacao deste comando abaixo no arquivo de bkp linha 141
	arquivo, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	dataHoje := time.Now().Format("02/01/2006" + " - " + "15:04")

	if err != nil {
		fmt.Println(err)
	}

	arquivo.WriteString(servico + " - ONLINE: " + strconv.FormatBool(status) + " - " + dataHoje + "\n")

	arquivo.Close()
}

func imprimeLogs() {
	fmt.Println("LISTAGEM DE LOGS ...")

	arquivo, err := ioutil.ReadFile("logs.txt") //Explicação no arquivo txt de bkp, na linha

	if err != nil {
		fmt.Println(err)
	}

	//imprime a listagem de logs
	fmt.Println(string(arquivo))
}
