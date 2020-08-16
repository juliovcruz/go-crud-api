# API CRUD With Golang

Este projeto foi criado para aprimorar meus conhecimentos em GoLang, se consiste em uma API simples que possui operações básicas CRUD.

## Tecnologias

- GoLang
- Protobuff
- gRPC
- MongoDB

## Como iniciar o projeto

Necessário possuir o compilador de GoLang instalado no computador

1. Clone o repositório
2. Inicie o servidor com o comando: "go run server/main.go"
3. Inicie o client com o comando: "go run client/main.go"

## Sobre o projeto

Foi utilizado a linguagem de programação Go,para persistência de dados foi utilizado o banco de dados não relacional MongoDB. O servidor utiliza o modelo gRPC para receber os comandos, e os dados e métodos são padronizados utilizando-se Protobuff.

A API possui 5 métodos básicos:

- CreateUser
- ReadUser
- UpdateUser
- DeleteUser
- ListUser