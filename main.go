package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

var (
	totalKeysChecked int64
	startTime        = time.Now()
	done             = make(chan struct{})
	ultimoEndereco   string
)

func carregarEnderecos(arquivo string) (map[string]bool, error) {
	enderecos := make(map[string]bool)
	file, err := os.Open(arquivo)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		endereco := scanner.Text()
		enderecos[endereco] = true
	}
	return enderecos, scanner.Err()
}

func gerarChaveAleatoria() []byte {
	privKey := make([]byte, 32)
	rand.Read(privKey)
	return privKey
}

func checarEnderecosAleatorios(enderecos map[string]bool, arquivoSaida *os.File, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-done:
			return
		default:
			privKeyBytes := gerarChaveAleatoria()
			privKey := secp256k1.PrivKeyFromBytes(privKeyBytes)
			pubKey := privKey.PubKey().SerializeCompressed()

			addressPubKey, err := btcutil.NewAddressPubKey(pubKey, &chaincfg.MainNetParams)
			if err != nil {
				fmt.Println("Erro ao criar endereço público:", err)
				return
			}

			addressSegWit, err := btcutil.NewAddressScriptHashFromHash(btcutil.Hash160(pubKey), &chaincfg.MainNetParams)
			if err != nil {
				fmt.Println("Erro ao criar endereço SegWit:", err)
				return
			}

			addressBech32, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(pubKey), &chaincfg.MainNetParams)
			if err != nil {
				fmt.Println("Erro ao criar endereço Bech32:", err)
				return
			}

			for _, addr := range []string{
				addressPubKey.EncodeAddress(),
				addressSegWit.EncodeAddress(),
				addressBech32.EncodeAddress(),
			} {
				ultimoEndereco = addr // Atualiza o último endereço checado
				if enderecos[addr] {
					fmt.Printf("\nChave Encontrada: %x\n", privKeyBytes)
					fmt.Printf("Endereço: %s\n", addr)
					arquivoSaida.WriteString(fmt.Sprintf("Endereço: %s - Chave Privada: %x\n", addr, privKeyBytes))
				}
			}
			totalKeysChecked++
		}
	}
}

func formatNumber(num int64) string {
	if num >= 1e12 {
		return fmt.Sprintf("%.2f TRILHÕES", float64(num)/1e12)
	} else if num >= 1e9 {
		return fmt.Sprintf("%.2f BILHÕES", float64(num)/1e9)
	} else if num >= 1e6 {
		return fmt.Sprintf("%.2f MILHÕES", float64(num)/1e6)
	} else if num >= 1e3 {
		return fmt.Sprintf("%.2f MIL", float64(num)/1e3)
	} else {
		return fmt.Sprintf("%d", num)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Uso: go run main.go <arquivo_enderecos> <arquivo_saida>")
		return
	}

	arquivoEnderecos := os.Args[1]
	enderecos, err := carregarEnderecos(arquivoEnderecos)
	if err != nil {
		fmt.Println("Erro ao carregar o arquivo de endereços:", err)
		return
	}

	arquivoSaida, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println("Erro ao criar o arquivo de saída:", err)
		return
	}
	defer arquivoSaida.Close()

	numCPUs := 15
	var wg sync.WaitGroup

	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go checarEnderecosAleatorios(enderecos, arquivoSaida, &wg)
	}

	// Goroutine para exibir estatísticas
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\033[H\033[J")
				fmt.Printf("      ADDRESS CHECKER 1.0 - Dev= ALAN ALVES 30/10/2024 BRASIL\n")
				fmt.Printf("      ------------------------------\n")
				fmt.Printf("    ENDEREÇOS VERIFICADOS: %15s\n", formatNumber(totalKeysChecked))
				fmt.Printf("  ENDEREÇOS P/S: %10.2f\n", float64(totalKeysChecked)/time.Since(startTime).Seconds())
				fmt.Printf("ÚLTIMO ENDEREÇO CHECADO: %s\n", ultimoEndereco)
				time.Sleep(2 * time.Second)
			}
		}
	}()

	wg.Wait()

	close(done)
	time.Sleep(1 * time.Second)
}

