package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/aseptimu/GophKeeper/internal/client"
	"golang.org/x/term"
)

func main() {
	var (
		serverAddr = flag.String("server", "127.0.0.1:8087", "Server address")
		help       = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	clientApp := client.NewClient(*serverAddr)

	fmt.Println("🔐 GophKeeper Client")
	fmt.Println("===================")

	if clientApp.IsAuthenticated() {
		fmt.Println("✅ You are already authenticated!")
		showAuthenticatedMenu(clientApp)
	} else {
		fmt.Println("❌ You are not authenticated.")
		showAuthMenu(clientApp)
	}
}

func showAuthMenu(clientApp *client.Client) {
	for {
		fmt.Println("\n📋 Authentication Menu:")
		fmt.Println("1. Login")
		fmt.Println("2. Register")
		fmt.Println("3. Exit")
		fmt.Print("Choose an option (1-3): ")

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			handleLogin(clientApp)
		case "2":
			handleRegister(clientApp)
		case "3":
			fmt.Println("👋 Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("❌ Invalid choice. Please try again.")
		}
	}
}

func showAuthenticatedMenu(clientApp *client.Client) {
	for {
		fmt.Println("\n🔒 Private Data Menu:")
		fmt.Println("1. View all data")
		fmt.Println("2. Add login/password")
		fmt.Println("3. Add bank card")
		fmt.Println("4. Add text data")
		fmt.Println("5. Add binary data (file)")
		fmt.Println("6. Logout")
		fmt.Println("7. Exit")
		fmt.Print("Choose an option (1-7): ")

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			handleViewData(clientApp)
		case "2":
			handleAddLoginPassword(clientApp)
		case "3":
			handleAddBankCard(clientApp)
		case "4":
			handleAddTextData(clientApp)
		case "5":
			handleAddBinaryData(clientApp)
		case "6":
			handleLogout(clientApp)
			return
		case "7":
			fmt.Println("👋 Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("❌ Invalid choice. Please try again.")
		}
	}
}

func handleLogin(clientApp *client.Client) {
	fmt.Print("Enter login: ")
	reader := bufio.NewReader(os.Stdin)
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)

	fmt.Print("Enter password: ")
	password, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	err := clientApp.Login(login, string(password))
	if err != nil {
		fmt.Printf("❌ Login failed: %v\n", err)
		return
	}

	fmt.Println("✅ Login successful!")
	showAuthenticatedMenu(clientApp)
}

func handleRegister(clientApp *client.Client) {
	fmt.Print("Enter login: ")
	reader := bufio.NewReader(os.Stdin)
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)

	fmt.Print("Enter password: ")
	password, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	err := clientApp.Register(login, string(password))
	if err != nil {
		fmt.Printf("❌ Registration failed: %v\n", err)
		return
	}

	fmt.Println("✅ Registration successful!")
	showAuthenticatedMenu(clientApp)
}

func handleViewData(clientApp *client.Client) {
	fmt.Println("\n📊 Your Private Data:")
	fmt.Println("====================")

	items, err := clientApp.GetDataItems()
	if err != nil {
		fmt.Printf("❌ Failed to get data: %v\n", err)
		return
	}

	if len(items) == 0 {
		fmt.Println("📝 No data found. Add some data first!")
		return
	}

	for i, item := range items {
		fmt.Printf("%d. [%s] %s\n", i+1, item.Type, item.Metadata)
		if item.Metadata != "" {
			fmt.Printf("   Metadata: %s\n", item.Metadata)
		}
		if item.Type == "binary" {
			// Для бинарных данных показываем дополнительную информацию
			var binaryData client.BinaryData
			if err := json.Unmarshal([]byte(item.Data), &binaryData); err == nil {
				fmt.Printf("   File: %s (%d bytes)\n", binaryData.FileName, binaryData.FileSize)
				fmt.Printf("   Type: %s\n", binaryData.ContentType)
			} else {
				fmt.Printf("   Data: Binary file data\n")
			}
		}
		fmt.Printf("   Created: %s\n", item.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}
}

func handleAddLoginPassword(clientApp *client.Client) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter login: ")
	login, _ := reader.ReadString('\n')
	login = strings.TrimSpace(login)

	fmt.Print("Enter password: ")
	password, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	fmt.Print("Enter metadata (optional): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)

	if login == "" || string(password) == "" {
		fmt.Println("❌ Login and password cannot be empty!")
		return
	}

	_, err := clientApp.CreateLoginPassword(login, string(password), metadata)
	if err != nil {
		fmt.Printf("❌ Failed to add login/password: %v\n", err)
		return
	}

	fmt.Println("✅ Login/password added successfully!")
}

func handleAddBankCard(clientApp *client.Client) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter card number: ")
	cardNumber, _ := reader.ReadString('\n')
	cardNumber = strings.TrimSpace(cardNumber)

	fmt.Print("Enter expiry date (MM/YY): ")
	expiryDate, _ := reader.ReadString('\n')
	expiryDate = strings.TrimSpace(expiryDate)

	fmt.Print("Enter CVV: ")
	cvv, _ := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()

	fmt.Print("Enter cardholder name: ")
	cardholder, _ := reader.ReadString('\n')
	cardholder = strings.TrimSpace(cardholder)

	fmt.Print("Enter metadata (optional): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)

	if cardNumber == "" || expiryDate == "" || string(cvv) == "" || cardholder == "" {
		fmt.Println("❌ All card fields are required!")
		return
	}

	_, err := clientApp.CreateBankCard(cardNumber, expiryDate, string(cvv), cardholder, metadata)
	if err != nil {
		fmt.Printf("❌ Failed to add bank card: %v\n", err)
		return
	}

	fmt.Println("✅ Bank card added successfully!")
}

func handleAddTextData(clientApp *client.Client) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter text data: ")
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)

	fmt.Print("Enter metadata (optional): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)

	if text == "" {
		fmt.Println("❌ Text data cannot be empty!")
		return
	}

	_, err := clientApp.CreateTextData(text, metadata)
	if err != nil {
		fmt.Printf("❌ Failed to add text data: %v\n", err)
		return
	}

	fmt.Println("✅ Text data added successfully!")
}

func handleAddBinaryData(clientApp *client.Client) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter file path: ")
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	if filePath == "" {
		fmt.Println("❌ File path cannot be empty!")
		return
	}

	// Читаем файл
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("❌ Failed to read file: %v\n", err)
		return
	}

	// Определяем MIME тип
	contentType := "application/octet-stream"
	if strings.HasSuffix(strings.ToLower(filePath), ".txt") {
		contentType = "text/plain"
	} else if strings.HasSuffix(strings.ToLower(filePath), ".pdf") {
		contentType = "application/pdf"
	} else if strings.HasSuffix(strings.ToLower(filePath), ".jpg") || strings.HasSuffix(strings.ToLower(filePath), ".jpeg") {
		contentType = "image/jpeg"
	} else if strings.HasSuffix(strings.ToLower(filePath), ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(strings.ToLower(filePath), ".zip") {
		contentType = "application/zip"
	}

	// Получаем имя файла
	fileName := filepath.Base(filePath)

	fmt.Print("Enter metadata (optional): ")
	metadata, _ := reader.ReadString('\n')
	metadata = strings.TrimSpace(metadata)

	if metadata == "" {
		metadata = fmt.Sprintf("File: %s (%d bytes)", fileName, len(fileData))
	}

	_, err = clientApp.CreateBinaryData(fileData, contentType, fileName, metadata)
	if err != nil {
		fmt.Printf("❌ Failed to add binary data: %v\n", err)
		return
	}

	fmt.Printf("✅ Binary data added successfully! (%d bytes)\n", len(fileData))
}

func handleLogout(clientApp *client.Client) {
	clientApp.Logout()
	fmt.Println("✅ Logged out successfully!")
	showAuthMenu(clientApp)
}

func showHelp() {
	fmt.Println("GophKeeper Interactive Client")
	fmt.Println("Usage: client [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -server string")
	fmt.Println("        Server address (default \"127.0.0.1:8087\")")
	fmt.Println("  -help")
	fmt.Println("        Show this help message")
	fmt.Println()
	fmt.Println("The client will start in interactive mode and guide you through")
	fmt.Println("authentication and data management.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  client")
	fmt.Println("  client -server localhost:8080")
}
