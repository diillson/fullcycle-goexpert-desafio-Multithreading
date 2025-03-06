Este código:

    Define estruturas para as respostas de ambas as APIs
    Usa goroutines para fazer requisições simultâneas
    Utiliza um canal para receber os resultados
    Implementa um timeout de 1 segundo usando context
    Mostra o resultado da API mais rápida
    Formata e exibe os dados do endereço no terminal

Para usar o programa:

    Execute com:

go run multithreading.go

O programa irá:

    Fazer requisições simultâneas para ambas as APIs
    Mostrar o resultado da API mais rápida
    Exibir um erro se ambas as APIs demorarem mais de 1 segundo
    Mostrar qual API forneceu o resultado

Você pode modificar o CEP na variável cep no início da função main para testar com diferentes CEPs.

O programa atende a todos os requisitos solicitados:

    Faz requisições simultâneas
    Mostra apenas o resultado mais rápido
    Exibe os dados do endereço e qual API respondeu
    Limita o tempo de resposta em 1 segundo
