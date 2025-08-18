//go:build linux

package connection

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/oLenador/mulltbost/internal/core/domain/entities"
)

type LinuxDNSExecutor struct{}

func NewDNSExecutor() *LinuxDNSExecutor {
	return &LinuxDNSExecutor{}
}

func (e *LinuxDNSExecutor) Execute(ctx context.Context, boosterID string) (*entities.BoosterResult, error) {
	// Backup das configurações atuais
	originalDNS, err := e.getCurrentDNSConfig()
	if err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: fmt.Sprintf("Falha ao obter configurações DNS atuais: %v", err),
		}, err
	}

	// DNS otimizados (Google DNS, Cloudflare, OpenDNS)
	optimizedDNS := []string{"8.8.8.8", "8.8.4.4", "1.1.1.1", "1.0.0.1", "208.67.222.222", "208.67.220.220"}

	// Aplicar nova configuração DNS
	if err := e.setDNSServers(optimizedDNS); err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: fmt.Sprintf("Falha ao configurar servidores DNS: %v", err),
		}, err
	}

	// Limpar cache DNS (systemd-resolved)
	if err := e.flushDNSCache(ctx); err != nil {
		// Não falhar se não conseguir limpar cache
		fmt.Printf("Aviso: Falha ao limpar cache DNS: %v\n", err)
	}

	// Otimizar configurações de rede via sysctl
	if err := e.optimizeNetworkSettings(ctx); err != nil {
		fmt.Printf("Aviso: Falha ao otimizar configurações de rede: %v\n", err)
	}

	return &entities.BoosterResult{
		Success: true,
		Message: "Configurações DNS otimizadas com sucesso para Linux. Cache DNS limpo e configurações de rede otimizadas.",
		BackupData: map[string]interface{}{
			"original_resolv_conf": originalDNS,
			"optimized_dns_servers": optimizedDNS,
			"backup_timestamp": time.Now().Unix(),
			"booster_id": boosterID,
			"platform": "linux",
		},
	}, nil
}

func (e *LinuxDNSExecutor) Revert(ctx context.Context, backupData entities.BackupData) (*entities.BoosterResult, error) {
	// Extrair configuração DNS original do backup
	originalConfigInterface, exists := backupData["original_resolv_conf"]
	if !exists {
		return &entities.BoosterResult{
			Success: false,
			Message: "Dados de backup não encontrados para restaurar configurações DNS",
		}, fmt.Errorf("backup data missing original_resolv_conf")
	}

	originalConfig, ok := originalConfigInterface.(string)
	if !ok {
		return &entities.BoosterResult{
			Success: false,
			Message: "Formato inválido dos dados de backup DNS",
		}, fmt.Errorf("invalid backup data format")
	}

	// Restaurar configuração original
	if err := e.restoreOriginalConfig(originalConfig); err != nil {
		return &entities.BoosterResult{
			Success: false,
			Message: fmt.Sprintf("Falha ao restaurar configurações DNS originais: %v", err),
		}, err
	}

	// Limpar cache DNS após restaurar
	if err := e.flushDNSCache(ctx); err != nil {
		fmt.Printf("Aviso: Falha ao limpar cache DNS após restaurar: %v\n", err)
	}

	// Restaurar configurações de rede padrão
	if err := e.restoreNetworkSettings(ctx); err != nil {
		fmt.Printf("Aviso: Falha ao restaurar configurações de rede: %v\n", err)
	}

	return &entities.BoosterResult{
		Success: true,
		Message: "Configurações DNS originais restauradas com sucesso",
	}, nil
}

func (e *LinuxDNSExecutor) Validate(ctx context.Context) error {
	// Verificar se estamos executando como root
	if os.Geteuid() != 0 {
		return fmt.Errorf("permissões de root necessárias para modificar configurações DNS")
	}

	// Verificar se o arquivo resolv.conf existe e é modificável
	resolvConfPath := "/etc/resolv.conf"
	if _, err := os.Stat(resolvConfPath); err != nil {
		return fmt.Errorf("arquivo /etc/resolv.conf não acessível: %v", err)
	}

	// Verificar se conseguimos escrever no arquivo
	file, err := os.OpenFile(resolvConfPath, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("sem permissão para escrever em /etc/resolv.conf: %v", err)
	}
	file.Close()

	// Verificar se comandos necessários estão disponíveis
	commands := []string{"systemctl", "sysctl"}
	for _, cmd := range commands {
		if _, err := exec.LookPath(cmd); err != nil {
			fmt.Printf("Aviso: comando %s não encontrado, algumas funcionalidades podem não estar disponíveis\n", cmd)
		}
	}

	return nil
}

func (e *LinuxDNSExecutor) CanExecute(ctx context.Context) bool {
	return e.Validate(ctx) == nil
}

// Métodos auxiliares

func (e *LinuxDNSExecutor) getCurrentDNSConfig() (string, error) {
	resolvConfPath := "/etc/resolv.conf"
	content, err := os.ReadFile(resolvConfPath)
	if err != nil {
		return "", fmt.Errorf("falha ao ler %s: %v", resolvConfPath, err)
	}
	
	return string(content), nil
}

func (e *LinuxDNSExecutor) setDNSServers(dnsServers []string) error {
	resolvConfPath := "/etc/resolv.conf"
	
	// Criar backup do arquivo atual
	backupPath := resolvConfPath + ".backup." + fmt.Sprintf("%d", time.Now().Unix())
	if err := e.copyFile(resolvConfPath, backupPath); err != nil {
		return fmt.Errorf("falha ao criar backup de resolv.conf: %v", err)
	}

	// Criar novo conteúdo do resolv.conf
	var content strings.Builder
	content.WriteString("# DNS configuration optimized by MultBoost\n")
	content.WriteString("# Backup available at: " + backupPath + "\n")
	content.WriteString("# Generated at: " + time.Now().Format(time.RFC3339) + "\n\n")
	
	// Adicionar servidores DNS otimizados
	for _, dns := range dnsServers {
		content.WriteString(fmt.Sprintf("nameserver %s\n", dns))
	}
	
	// Adicionar configurações adicionais para melhor performance
	content.WriteString("\n# Performance optimizations\n")
	content.WriteString("options timeout:1\n")
	content.WriteString("options attempts:3\n")
	content.WriteString("options rotate\n")

	// Escrever novo arquivo
	if err := os.WriteFile(resolvConfPath, []byte(content.String()), 0644); err != nil {
		return fmt.Errorf("falha ao escrever novo resolv.conf: %v", err)
	}

	return nil
}

func (e *LinuxDNSExecutor) restoreOriginalConfig(originalConfig string) error {
	resolvConfPath := "/etc/resolv.conf"
	
	if err := os.WriteFile(resolvConfPath, []byte(originalConfig), 0644); err != nil {
		return fmt.Errorf("falha ao restaurar resolv.conf original: %v", err)
	}
	
	return nil
}

func (e *LinuxDNSExecutor) flushDNSCache(ctx context.Context) error {
	// Tentar diferentes métodos de limpeza de cache DNS

	// systemd-resolved
	if err := exec.CommandContext(ctx, "systemctl", "is-active", "systemd-resolved").Run(); err == nil {
		cmd := exec.CommandContext(ctx, "systemd-resolve", "--flush-caches")
		if err := cmd.Run(); err != nil {
			// Tentar comando alternativo
			cmd = exec.CommandContext(ctx, "resolvectl", "flush-caches")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("falha ao limpar cache systemd-resolved: %v", err)
			}
		}
		return nil
	}

	// nscd (Name Service Cache Daemon)
	if _, err := exec.LookPath("nscd"); err == nil {
		cmd := exec.CommandContext(ctx, "nscd", "-i", "hosts")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("falha ao limpar cache nscd: %v", err)
		}
		return nil
	}

	// dnsmasq
	if err := exec.CommandContext(ctx, "systemctl", "is-active", "dnsmasq").Run(); err == nil {
		cmd := exec.CommandContext(ctx, "systemctl", "restart", "dnsmasq")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("falha ao reiniciar dnsmasq: %v", err)
		}
		return nil
	}

	return fmt.Errorf("nenhum serviço de cache DNS encontrado para limpar")
}

func (e *LinuxDNSExecutor) optimizeNetworkSettings(ctx context.Context) error {
	// Configurações de rede otimizadas
	settings := map[string]string{
		"net.core.rmem_default":     "31457280",
		"net.core.rmem_max":         "67108864",
		"net.core.wmem_default":     "31457280", 
		"net.core.wmem_max":         "67108864",
		"net.core.somaxconn":        "65535",
		"net.core.netdev_max_backlog": "5000",
		"net.ipv4.tcp_congestion_control": "bbr",
		"net.ipv4.tcp_rmem":         "4096 31457280 67108864",
		"net.ipv4.tcp_wmem":         "4096 31457280 67108864",
		"net.ipv4.tcp_fastopen":     "3",
		"net.ipv4.tcp_slow_start_after_idle": "0",
	}

	for key, value := range settings {
		cmd := exec.CommandContext(ctx, "sysctl", "-w", fmt.Sprintf("%s=%s", key, value))
		if err := cmd.Run(); err != nil {
			fmt.Printf("Aviso: Falha ao configurar %s: %v\n", key, err)
		}
	}

	return nil
}

func (e *LinuxDNSExecutor) restoreNetworkSettings(ctx context.Context) error {
	// Restaurar configurações de rede padrão
	defaultSettings := map[string]string{
		"net.core.rmem_default":     "212992",
		"net.core.rmem_max":         "212992", 
		"net.core.wmem_default":     "212992",
		"net.core.wmem_max":         "212992",
		"net.core.somaxconn":        "128",
		"net.core.netdev_max_backlog": "1000",
		"net.ipv4.tcp_congestion_control": "cubic",
		"net.ipv4.tcp_rmem":         "4096 87380 6291456",
		"net.ipv4.tcp_wmem":         "4096 16384 4194304", 
		"net.ipv4.tcp_fastopen":     "1",
		"net.ipv4.tcp_slow_start_after_idle": "1",
	}

	for key, value := range defaultSettings {
		cmd := exec.CommandContext(ctx, "sysctl", "-w", fmt.Sprintf("%s=%s", key, value))
		if err := cmd.Run(); err != nil {
			fmt.Printf("Aviso: Falha ao restaurar %s: %v\n", key, err)
		}
	}

	return nil
}

func (e *LinuxDNSExecutor) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return err
	}

	buffer := make([]byte, sourceInfo.Size())
	_, err = sourceFile.Read(buffer)
	if err != nil {
		return err
	}

	_, err = destFile.Write(buffer)
	return err
}