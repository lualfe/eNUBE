DROP TABLE IF EXISTS dados;

CREATE TABLE dados
(
     id SERIAL PRIMARY KEY,
     nome VARCHAR(50),
     inicio DATE,
     fim DATE,
     quantidade INT,
     unidade VARCHAR(50),
     preco NUMERIC(15,2)
)