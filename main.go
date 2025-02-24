package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

// Entry representa uma única entrada na tabela Markdown.
type Entry struct {
	Canonico    string // Caminho canônico do atributo na estrutura JSON
	Descricao   string // Descrição do atributo
	Nome        string // Nome do atributo
	Obrigatorio string // Indica se o atributo é obrigatório
	Tipo        string // Tipo do atributo
}

func main() {
	// Define flags de linha de comando para arquivos de entrada e saída
	inputFile := flag.String("j", "", "Nome do arquivo JSON de entrada")
	outputFile := flag.String("o", "", "Nome do arquivo Markdown de saída")
	flag.Parse()

	if *inputFile != "" {
		data, err := readFile(*inputFile)
		if err != nil {
			log.Fatal(err)
		}
		result, err := generateMarkdown(data)
		if err != nil {
			log.Fatal(err)
		}
		if *outputFile != "" {
			err = os.WriteFile(*outputFile, []byte(result), 0644)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println(result)
		}
	} else {
		http.HandleFunc("/", homeHandler)
		http.HandleFunc("/convert", convertHandler)
		log.Println("Servidor iniciado na porta 3000")
		log.Fatal(http.ListenAndServe(":3000", nil))
	}
}

func renderTemplate(w http.ResponseWriter, jsonInput, markdownResult string) {
	tmpl := `
    <!DOCTYPE html>
    <html>
    <head>
        <title>JSON to Markdown</title>
        <style>
            textarea {
                font-family: monospace;
            }
        </style>
        <script>
            function handleFileSelect(evt) {
                var file = evt.target.files[0];
                if (file) {
                    var reader = new FileReader();
                    reader.onload = function(e) {
                        document.getElementById('jsonInput').value = e.target.result;
                    };
                    reader.readAsText(file);
                }
            }
        </script>
    </head>
    <body>
        <h1>Conversor de JSON para Markdown</h1>
        <form action="/convert" method="post">
            <label for="jsonFile">Arquivo JSON:</label>
            <input type="file" id="jsonFile" name="jsonFile" onchange="handleFileSelect(event)"><br><br>
            <label for="jsonInput">Ou cole seu JSON aqui:</label><br>
            <textarea id="jsonInput" name="jsonInput" rows="10" cols="100">{{.JSONInput}}</textarea><br><br>
            <input type="submit" value="Converter"><br><br>
            <label for="markdownResult">Resultado em Markdown:</label><br>
            <textarea id="markdownResult" name="markdownResult" rows="20" cols="100">{{.Markdown}}</textarea><br><br>
        </form>
    </body>
    </html>
    `
	t := template.Must(template.New("result").Parse(tmpl))
	t.Execute(w, map[string]string{"JSONInput": jsonInput, "Markdown": markdownResult})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "", "")
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	jsonInput := r.FormValue("jsonInput")
	if jsonInput == "" {
		http.Error(w, "Nenhum JSON fornecido", http.StatusBadRequest)
		return
	}

	var data interface{}
	if err := json.Unmarshal([]byte(jsonInput), &data); err != nil {
		http.Error(w, "Erro ao fazer unmarshal do JSON", http.StatusInternalServerError)
		return
	}

	// Gera o Markdown
	markdownResult, err := generateMarkdown([]byte(jsonInput))
	if err != nil {
		http.Error(w, "Erro ao gerar Markdown", http.StatusInternalServerError)
		return
	}

	// Renderiza o template com o resultado
	renderTemplate(w, jsonInput, markdownResult)
}

// readFile lê o conteúdo de um arquivo e retorna os bytes lidos.
func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

// generateMarkdown gera o Markdown a partir dos dados JSON.
func generateMarkdown(data []byte) (string, error) {
	var markdownResult strings.Builder
	markdownResult.WriteString("# Tabela de Atributos\n\n")
	markdownResult.WriteString("| <span style=\"color:#1F4E79;\">Nome | <span style=\"color:#1F4E79;\">Descrição | <span style=\"color:#1F4E79;\">Tipo | <span style=\"color:#1F4E79;\">Obrigatório | <span style=\"color:#1F4E79;\">Nome Canônico |\n")
	markdownResult.WriteString("| :--- | :-------- | :--- | :----------: | :------------ |\n")

	entries := []Entry{}
	processJSON(gjson.ParseBytes(data), "", &entries, "")

	for _, entry := range entries {
		line := fmt.Sprintf("| %s | %s | %s | %s | %s |\n", entry.Nome, entry.Descricao, entry.Tipo, entry.Obrigatorio, entry.Canonico)
		line = strings.ReplaceAll(line, "{}{}", "{}")
		line = strings.ReplaceAll(line, "->->", "->")
		markdownResult.WriteString(line)
	}

	return markdownResult.String(), nil
}

// processJSON processa recursivamente os dados JSON para gerar entradas para a tabela Markdown.
func processJSON(result gjson.Result, path string, entries *[]Entry, parent string) {
	if result.IsObject() {
		if parent != "" && !strings.Contains(path, "[]") {
			// Verifica se o nome do atributo termina com "?"
			obrigatorio := "sim"
			if strings.HasSuffix(parent, "?") {
				obrigatorio = "não"
				parent = strings.TrimSuffix(parent, "?")
			}

			entry := Entry{
				Nome:        parent,
				Descricao:   "-",
				Obrigatorio: obrigatorio,
				Tipo:        "objeto",
				Canonico:    "{}->" + path + "{}",
			}
			*entries = append(*entries, entry)
		}
		result.ForEach(func(key, value gjson.Result) bool {
			newPath := key.String()
			newPath = strings.TrimSuffix(newPath, "?")
			if path != "" {
				newPath = path + "->" + newPath
			}
			processJSON(value, newPath, entries, key.String())
			return true
		})
	} else if result.IsArray() {
		arrayPath := path + "[]"
		entry := Entry{
			Nome:        parent,
			Descricao:   "-",
			Obrigatorio: "não",
			Tipo:        "array",
			Canonico:    "{}->" + arrayPath + "->{}",
		}
		*entries = append(*entries, entry)
		if len(result.Array()) > 0 {
			processJSON(result.Array()[0], arrayPath+"->{}", entries, "")
		}
	} else {
		// Verifica se o nome do atributo termina com "?"
		obrigatorio := "sim"
		if strings.HasSuffix(parent, "?") {
			obrigatorio = "não"
			parent = strings.TrimSuffix(parent, "?")
		}

		entry := Entry{
			Nome:        parent,
			Descricao:   "-",
			Obrigatorio: obrigatorio,
			Tipo:        getType(result.Value()),
			Canonico:    "{}->" + path,
		}
		*entries = append(*entries, entry)
	}
}

// getType retorna o tipo de um valor.
func getType(val interface{}) string {
	switch v := val.(type) {
	case string:
		return "string"
	case float64:
		if v == float64(int(v)) {
			return "inteiro"
		}
		return "decimal"
	case map[string]interface{}:
		return "objeto"
	case []interface{}:
		return "array"
	case bool:
		return "booleano"
	default:
		return "desconhecido"
	}
}
