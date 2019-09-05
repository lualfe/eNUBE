# Leitura de Planilhas

Essa aplicação lê os dados de uma planilha ou .txt e retorna um resumo dos dados.

1. Para testar a API, é necessário criar uma database no psql shell 
   ```CREATE DATABASE nome_do_banco``` 
   e executar o arquivo setup.sql na CLI do PostgreSQL.
2. No arquivo main.go, ir na função define variables. Lá, onde há aspas vazias, deve-se colocar os dados de acordo com o banco a ser utilizado.
3. Após atualizar as variáveis e salvar o arquivo, basta acessar a pasta root do arquivo main.go no terminal e executar o comando go build.
4. Com o binário gerado rodando, basta acessar o [postman](https://documenter.getpostman.com/view/8679941/SVfUsmwD) e executar o POST seguido do GET.