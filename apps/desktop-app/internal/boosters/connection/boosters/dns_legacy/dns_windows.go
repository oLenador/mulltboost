//go:build windows

package connection

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type WindowsDNSExecutor struct{}

func NewDNSExecutor() *WindowsDNSExecutor {
	return &WindowsDNSExecutor{}
}

func (e *WindowsDNSExecutor) Execute(ctx context.Context, boosterID string) (*entities.BoosterResult, error) {
	// Backup das configurações atuais
	originalDNS, err := e.getCurrentDNSServers(ctx)
	if err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: fmt.Sprintf("Falha ao obter configurações DNS atuais: %v", err),
		}, err
	}

	// Configurar DNS otimizado (Google DNS e Cloudflare)
	optimizedDNS := []string{"8.8.8.8", "8.8.4.4", "1.1.1.1", "1.0.0.1"}
	
	// Aplicar nova configuração DNS
	if err := e.setDNSServers(ctx, optimizedDNS); err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: fmt.Sprintf("Falha ao configurar servidores DNS: %v", err),
		}, err
	}

	// Limpar cache DNS
	if err := e.flushDNSCache(ctx); err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: fmt.Sprintf("Falha ao limpar cache DNS: %v", err),
		}, err
	}

	// Renovar configuração de rede
	if err := e.renewNetworkConfig(ctx); err != nil {
		// Não falhar completamente se não conseguir renovar
		fmt.Printf("Aviso: Falha ao renovar configuração de rede: %v\n", err)
	}

	return &entities.BoosterResult{
		Success: true,
		Message: "Configurações DNS otimizadas com sucesso. Cache DNS limpo e configuração de rede renovada.",
		BackupData: map[string]interface{}{
			"original_dns_servers": originalDNS,
			"optimized_dns_servers": optimizedDNS,
			"backup_timestamp": time.Now().Unix(),
			"booster_id": boosterID,
			"platform": "windows",
		},
	}, nil
}

func (e *WindowsDNSExecutor) Revert(ctx context.Context, backupData entities.BackupData) (*entities.BoosterResult, error) {
	// Extrair servidores DNS originais do backup
	originalDNSInterface, exists := backupData["original_dns_servers"]
	if !exists {
		return &entities.BoosterResult{
			Success: false,
			Message: "Dados de backup não encontrados para restaurar configurações DNS",
		}, fmt.Errorf("backup data missing original_dns_servers")
	}

	originalDNS, ok := originalDNSInterface.([]string)
	if !ok {
		return &entities.BoosterResult{
			Success: false,
			Message: "Formato inválido dos dados de backup DNS",
		}, fmt.Errorf("invalid backup data format")
	}

	// Restaurar configuração DNS original
	if len(originalDNS) > 0 {
		if err := e.setDNSServers(ctx, originalDNS); err != nil {
			return &entities.BoosterResult{
				Success: false,
				Message: fmt.Sprintf("Falha ao restaurar servidores DNS originais: %v", err),
			}, err
		}
	} else {
		// Se não havia DNS configurados, configurar para automático
		if err := e.setAutomaticDNS(ctx); err != nil {
			return &entities.BoosterResult{
				Success: false,
				Message: fmt.Sprintf("Falha ao configurar DNS automático: %v", err),
			}, err
		}
	}

	// Limpar cache DNS após restaurar
	if err := e.flushDNSCache(ctx); err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: fmt.Sprintf("Falha ao limpar cache DNS após restaurar: %v", err),
		}, err
	}

	return &entities.BoosterResult{
		Success: true,
		Message: "Configurações DNS originais restauradas com sucesso",
	}, nil
}

func (e *WindowsDNSExecutor) Validate(ctx context.Context) error {
	// Verificar se os comandos necessários estão disponíveis
	commands := []string{"netsh", "ipconfig"}
	
	for _, cmd := range commands {
		if _, err := exec.LookPath(cmd); err != nil {
			return fmt.Errorf("comando necessário não encontrado: %s", cmd)
		}
	}

	// Verificar se temos permissões administrativas
	cmd := exec.CommandContext(ctx, "net", "session")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("permissões administrativas necessárias para modificar configurações DNS")
	}

	// Verificar se conseguimos listar interfaces de rede
	if _, err := e.getActiveNetworkInterface(ctx); err != nil {
		return fmt.Errorf("não foi possível obter interface de rede ativa: %v", err)
	}

	return nil
}

func (e *WindowsDNSExecutor) CanExecute(ctx context.Context) bool {
	return e.Validate(ctx) == nil
}

// Métodos auxiliares

func (e *WindowsDNSExecutor) getCurrentDNSServers(ctx context.Context) ([]string, error) {
	interfaceName, err := e.getActiveNetworkInterface(ctx)
	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, "netsh", "interface", "ipv4", "show", "dns", interfaceName)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var dnsServers []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "DNS Servers") {
			continue
		}
		if strings.Contains(line, ".") && len(strings.Split(line, ".")) == 4 {
			// Parece um endereço IP
			ip := strings.TrimSpace(line)
			if e.isValidIP(ip) {
				dnsServers = append(dnsServers, ip)
			}
		}
	}

	return dnsServers, nil
}

func (e *WindowsDNSExecutor) setDNSServers(ctx context.Context, dnsServers []string) error {
	interfaceName, err := e.getActiveNetworkInterface(ctx)
	if err != nil {
		return err
	}

	// Primeiro, limpar configurações DNS existentes
	cmd := exec.CommandContext(ctx, "netsh", "interface", "ipv4", "set", "dns", interfaceName, "dhcp")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("falha ao limpar configurações DNS: %v", err)
	}

	// Configurar servidores DNS
	for i, dns := range dnsServers {
		var cmd *exec.Cmd
		if i == 0 {
			// Primeiro servidor DNS (primário)
			cmd = exec.CommandContext(ctx, "netsh", "interface", "ipv4", "set", "dns", interfaceName, "static", dns, "primary")
		} else {
			// Servidores DNS adicionais
			cmd = exec.CommandContext(ctx, "netsh", "interface", "ipv4", "add", "dns", interfaceName, dns)
		}
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("falha ao configurar servidor DNS %s: %v", dns, err)
		}
	}

	return nil
}

func (e *WindowsDNSExecutor) setAutomaticDNS(ctx context.Context) error {
	interfaceName, err := e.getActiveNetworkInterface(ctx)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "netsh", "interface", "ipv4", "set", "dns", interfaceName, "dhcp")
	return cmd.Run()
}

func (e *WindowsDNSExecutor) flushDNSCache(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "ipconfig", "/flushdns")
	return cmd.Run()
}

func (e *WindowsDNSExecutor) renewNetworkConfig(ctx context.Context) error {
	// Release e renew da configuração IP
	releaseCmd := exec.CommandContext(ctx, "ipconfig", "/release")
	if err := releaseCmd.Run(); err != nil {
		return fmt.Errorf("falha ao liberar configuração IP: %v", err)
	}

	renewCmd := exec.CommandContext(ctx, "ipconfig", "/renew")
	if err := renewCmd.Run(); err != nil {
		return fmt.Errorf("falha ao renovar configuração IP: %v", err)
	}

	return nil
}

func (e *WindowsDNSExecutor) getActiveNetworkInterface(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "netsh", "interface", "show", "interface")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Connected") && strings.Contains(line, "Dedicated") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				// O nome da interface geralmente é a última parte
				return strings.Join(parts[3:], " "), nil
			}
		}
	}

	// Fallback para interface padrão
	return "Local Area Connection", nil
}

func (e *WindowsDNSExecutor) isValidIP(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	
	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
		}
	}
	
	return true
}