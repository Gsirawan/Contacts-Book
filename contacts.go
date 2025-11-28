package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const displayMenu = "(1). Add Contact \n(2). List Contacts \n(3). Search\n(4). Exit"

type Contact struct {
	Name   string
	Email  string
	Mobile string
}

// create new contacts
func addContact() {
	var name string
	var email string
	var mobile string
	var err error
	var file *os.File
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	var zeroCheck bool

	fmt.Println("Adding New Contact")

	// name input
	fmt.Println("Enter the new contact name:")
	fmt.Println("---------------------------")
	reader := bufio.NewReader(os.Stdin)
	name, err = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	upperC := cases.Title(language.English)
	name = upperC.String(name)
	if err != nil {
		log.Fatalf("Error reading input %v:\n %v:\n", name, err)
		return
	}

	// email input
	fmt.Println("Enter the new contact email:")
	fmt.Println("---------------------------")
	reader = bufio.NewReader(os.Stdin)
	for {
		email, err = reader.ReadString('\n')
		email = strings.TrimSpace(email)
		match, _ := regexp.MatchString(pattern, email)
		if !match {
			fmt.Println("Please Enter a valid email address!")
			fmt.Println("--------------------------------")
		} else {
			break
		}
	}
	if err != nil {
		log.Fatalf("Error reading input %v:\n\n%v:\n", email, err)
		return
	}

	// mobile input
	fmt.Println("Enter the new contact mobile:")
	fmt.Println("---------------------------")
	reader = bufio.NewReader(os.Stdin)
	for {
		mobile, err = reader.ReadString('\n')
		mobile = strings.TrimSpace(mobile)
		zeroCheck = strings.HasPrefix(mobile, "0")
		if len(mobile) != 10 || !zeroCheck {
			fmt.Println("Please Enter a valid mobile number, Start with '0'")
			fmt.Println("----------------------------------------------------")
		} else {
			mobile = strings.Replace(mobile, "0", "+971", 1)
			break
		}
	}
	if err != nil {
		log.Fatalf("Error reading input %v:\n\n%v:\n", mobile, err)
		return
	}
	// adding new contact
	newContact := Contact{
		Name:   name,
		Email:  email,
		Mobile: mobile,
	}

	// creare the file
	file, err = os.OpenFile("contacts.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalf("Error opening file %v\n:", err)
		return
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("Error closing file %v\n:", err)
		}
	}()

	// write on the file
	_, err = fmt.Fprintf(file, "%s,%s,%s\n", newContact.Name, newContact.Email, newContact.Mobile)
	if err != nil {
		log.Fatalf("Error writing to file %v\n:", err)
		return
	}
	fmt.Println("Successfully saved input")
	fmt.Printf("┃Name: %s\n┃Email: %s\n┃Mobile: %s\n", name, email, mobile)
	fmt.Println("----------------")
}

// List all the contents
func listContact() {
	var err error
	var file *os.File

	file, err = os.OpenFile("contacts.txt", os.O_RDONLY, 0o644)
	if err != nil {
		log.Fatalf("Error Opening the file %v\n", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("Error closing file %v\n:", err)
		}
	}()

	fmt.Println("--- List of Contents ---")
	scanner := bufio.NewScanner(file)
	fmt.Println("┃        Name        ┃        Email        ┃        Mobile        ┃")
	fmt.Println("=====================================================================")
	for scanner.Scan() {
		line := scanner.Text()
		part := strings.Split(line, ",")
		contact := Contact{
			Name:   part[0],
			Email:  part[1],
			Mobile: part[2],
		}
		fmt.Printf("|%-20s|%-21s|%-20s", contact.Name, contact.Email, contact.Mobile)
		fmt.Printf("\n")
	}
	fmt.Println("---------------------------------------")
}

func countContact() {
	var err error
	var file *os.File
	counter := 0

	file, err = os.OpenFile("contacts.txt", os.O_RDONLY, 0o644)
	if err != nil {
		log.Fatalf("Error Opening the file %v\n", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatalf("Error closing file %v\n:", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		counter++
	}
	fmt.Printf("Contact available: %d\n", counter)
	fmt.Println("---------------------------------------")
}

func search() {
	fmt.Println("Search for contact:")
	fmt.Println("------------------")
}

func main() {
	var choice string
	for {
		fmt.Println(displayMenu)
		fmt.Println("=====================")
		fmt.Println("Please choice from the list:")
		fmt.Println("---------------------")
		_, err := fmt.Scanf("%s", &choice)
		if err != nil {
			log.Fatalf("Error reading input %v\n:", err)
			return
		}
		if choice == "1" {
			addContact()
		} else if choice == "2" {
			listContact()
			countContact()
		} else if choice == "3" {
			search()
		} else if choice == "4" {
			fmt.Println("Goodbye!")
			break
		} else {
			fmt.Println("You can only choose from the List")
			// fmt.Scanln(&choice)
		}
	}
}
