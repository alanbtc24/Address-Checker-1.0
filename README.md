# Address-Checker-1.0
Geração e checagem de Endereços Bitcoin

Lógica do Programa

O programa implementa um verificador de endereços Bitcoin que realiza as seguintes operações:

    Carregar Endereços: Lê endereços Bitcoin de um arquivo de entrada e armazena-os em um mapa para fácil verificação.

    Gerar Chaves Aleatórias: Gera chaves privadas aleatórias.

    Converter Chaves em Endereços: Para cada chave gerada, o programa:
        Converte a chave privada em uma chave pública.
        Gera três tipos de endereços Bitcoin a partir da chave pública:
            Endereço legado (prefixo "1")
            Endereço SegWit (prefixo "3")
            Endereço Bech32 (prefixo "bc1")

    Verificar Endereços: Compara os endereços gerados com os endereços carregados do arquivo de entrada. Se encontrar uma correspondência, imprime a chave privada e o endereço, e os grava em um arquivo de saída.

    Execução Contínua: O programa continua a executar indefinidamente, gerando e verificando endereços até que o usuário decida parar.

    Relatório de Progresso: Durante a execução, o programa exibe estatísticas em tempo real, incluindo o número total de endereços verificados e a taxa de endereços verificados por segundo.


Comando para Executar

Para executar o programa, utilize o seguinte comando no terminal:

    
   go run main.go Bitcoin_addresses_LATEST.txt endereços.txt


Este comando irá executar o programa, verificando os endereços no arquivo Bitcoin_addresses_LATEST.txt e salvando os resultados encontrados no arquivo endereços.txt.
