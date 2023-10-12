# Api em golang que realiza a leitura de um pdf e retorna a linha digitavel e/ou a código de barras

Api construida em golang, usando echo v4 para criar um servidor web, tratando erros e retornando o resultado em json.
Será implementado novas melhorias e novas funcionalidades,
sendo atualmente possível ler a linha digitável que está em formato texto, e como fallback caso não encontre ou o boleto seja uma imagem,
irá sempre buscar por um código de barras, caso não encontre nenhum dos dois, retornará um erro.

## Funcionalidades

- Retorna a linha digitável sem formatação.
- Retorna o código de barras sem formatação.
- Retorna uma lista do que foi identificado.
- Leitura de pdf protegido por senha, basta informar o parâmetro password, passando a senha do pdf.


## Como usar com docker (recomendado)

- Faça o download do projeto
- Execute o comando 'docker-compose up -d'
- Faça uma requisição para o endereço `http://localhost:8080/upload`,
passando o parametro file com o arquivo pdf que deseja ler, caso o pdf esteja protegido por senha, passe o parametro password com a senha do pdf.

---
**Observação:**
caso queira rodar sem docker, basta seguir o passo a passo do arquivo Dockerfile, realizando as devidas instalações dos pacotes necessários. 
(poppler-utils, libzbar-dev) o nome do pacote pode variar de acordo com o sistema operacional.
---

## Exemplo de retorno json

```json
{
  "typeableLine": "23793381286008301352856000063307789840000150000",
  "barCode": "23797898400001500003381260083013525600006330",
  "findTypes": [
    "typeableLine",
    "barCode"
  ]
}
```

## Relacionado

[boleto-pdf-read
](https://github.com/fabioalmeida132/boleto-pdf-read) - Versão do projeto desenvolvido em node.


## Licença

[MIT](https://choosealicense.com/licenses/mit/)

