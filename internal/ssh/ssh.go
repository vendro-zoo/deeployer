package ssh

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

type Client struct {
	DryRun  bool
	Verbose bool
}

func New(dryRun, verbose bool) *Client {
	return &Client{
		DryRun:  dryRun,
		Verbose: verbose,
	}
}

func (c *Client) ExecuteCommands(host, user string, commands []string) error {
	if len(commands) == 0 {
		return nil
	}

	if c.DryRun {
		for _, cmd := range commands {
			fmt.Printf("Would execute on %s@%s: %s\n", user, host, cmd)
		}
		return nil
	}

	client, err := c.connect(host, user)
	if err != nil {
		return fmt.Errorf("failed to connect to %s@%s: %w", user, host, err)
	}
	defer client.Close()

	for _, command := range commands {
		if err := c.executeCommand(client, command); err != nil {
			return fmt.Errorf("command failed on %s@%s: %s: %w", user, host, command, err)
		}
	}

	return nil
}

func (c *Client) connect(host, user string) (*ssh.Client, error) {
	config, err := c.getSSHConfig(user)
	if err != nil {
		return nil, err
	}

	address := host + ":22"
	if strings.Contains(host, ":") {
		address = host
	}

	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) getSSHConfig(user string) (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if hostKeyCallback, err := c.getHostKeyCallback(); err == nil {
		config.HostKeyCallback = hostKeyCallback
	}

	authMethods, err := c.getAuthMethods()
	if err != nil {
		return nil, err
	}

	config.Auth = authMethods
	return config, nil
}

func (c *Client) getHostKeyCallback() (ssh.HostKeyCallback, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	knownHostsPath := filepath.Join(homeDir, ".ssh", "known_hosts")
	if _, err := os.Stat(knownHostsPath); os.IsNotExist(err) {
		return ssh.InsecureIgnoreHostKey(), nil
	}

	return knownhosts.New(knownHostsPath)
}

func (c *Client) getAuthMethods() ([]ssh.AuthMethod, error) {
	var authMethods []ssh.AuthMethod

	if agentAuth := c.getSSHAgent(); agentAuth != nil {
		authMethods = append(authMethods, agentAuth)
	}

	keyAuth, err := c.getPublicKeyAuth()
	if err == nil {
		authMethods = append(authMethods, keyAuth)
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication methods available")
	}

	return authMethods, nil
}

func (c *Client) getSSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func (c *Client) getPublicKeyAuth() (ssh.AuthMethod, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	keyPaths := []string{
		filepath.Join(homeDir, ".ssh", "id_rsa"),
		filepath.Join(homeDir, ".ssh", "id_ed25519"),
		filepath.Join(homeDir, ".ssh", "id_ecdsa"),
	}

	for _, keyPath := range keyPaths {
		if key, err := c.loadPrivateKey(keyPath); err == nil {
			return ssh.PublicKeys(key), nil
		}
	}

	return nil, fmt.Errorf("no private keys found")
}

func (c *Client) loadPrivateKey(path string) (ssh.Signer, error) {
	key, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return signer, nil
}

func (c *Client) executeCommand(client *ssh.Client, command string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if c.Verbose {
		fmt.Printf("Executing remote command: %s\n", command)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	return session.Run(command)
}