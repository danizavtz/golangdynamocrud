package model

type Usuario struct {
	Identificador string `json:"identifier"`
	Detalhes Detalhe `json:"detalhe"`
}

type Detalhe struct {
	Idade int `json:"idade" binding:"required"`
	Nome  string `json:"nome" binding:"required"`
	Profissao string `json:"profissao" binding:"required"`
}