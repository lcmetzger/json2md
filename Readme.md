# JSON to Markdown Converter

Este projeto é um conversor de JSON para Markdown que pode ser executado tanto como uma aplicação de linha de comando quanto como um servidor web.

Dado um arquivo JSON de entreda, será entregue uma tabela em Markdown com o mapeamento dos atributos e tipos, juntamente com uma estrutura canônica.

## Como Executar

### Linha de Comando

Para executar o conversor via linha de comando, use os seguintes comandos:

1. Compile o projeto:

```sh
  go build -o json2md main.go
```

1. Execute o binário gerado, especificando o arquivo JSON de entrada e o arquivo Markdown de saída:

```sh
  ./json2md -j teste.json -o resultado.md
```

### Servidor Web

Para executar o conversor como um servidor web, use o seguinte comando:

  ```sh
    ./json2md
  ```

## Como Executar no Windows

### Linha de Comando no Windows

Para executar o conversor via linha de comando no Windows, use os seguintes comandos:

1. Compile o projeto:

```sh
  go build -o json2md.exe main.go
```

1. Execute o binário gerado, especificando o arquivo JSON de entrada e o arquivo Markdown de saída:

```sh
  .\json2md.exe -j teste.json -o resultado.md
```

### Servidor Web no Windows

Para executar o conversor como um servidor web no Windows, use o seguinte comando:

```sh
  .\json2md
```
