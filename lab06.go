package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/godror/godror"
)

func transferMoney(dbA *sql.DB, dbB *sql.DB, amount int) error {
	// Start a distributed transaction for database A
	txA, err := dbA.Begin()
	if err != nil {
		return err
	}
	defer txA.Rollback()

	// Start a distributed transaction for database B
	txB, err := dbB.Begin()
	if err != nil {
		return err
	}
	defer txB.Rollback()

	// Lock customerA in database A
	_, err = txA.Exec("SELECT * FROM customerA WHERE id = 1 FOR UPDATE")
	if err != nil {
		return err
	}

	// Lock customerB in database B
	_, err = txB.Exec("SELECT * FROM customerB WHERE id = 2 FOR UPDATE")
	if err != nil {
		return err
	}

	// Get customerA's current balance in database A
	var customerABalance int
	err = txA.QueryRow("SELECT balance FROM customerA WHERE id = 1").Scan(&customerABalance)
	if err != nil {
		return err
	}

	// Update customerA's balance in database A
	_, err = txA.Exec("UPDATE customerA SET balance = balance - :1 WHERE id = 1", amount)
	if err != nil {
		return err
	}

	// Get customerB's current balance in database B
	var customerBBalance int
	err = txB.QueryRow("SELECT balance FROM customerB WHERE id = 2").Scan(&customerBBalance)
	if err != nil {
		return err
	}

	// Update customerB's balance in database B
	_, err = txB.Exec("UPDATE customerB SET balance = balance + :1 WHERE id = 2", amount)
	if err != nil {
		return err
	}

	// Commit the distributed transaction for database A
	if err := txA.Commit(); err != nil {
		return err
	}

	// Commit the distributed transaction for database B
	if err := txB.Commit(); err != nil {
		return err
	}

	return nil
}

func main() {
	// Connection details for Oracle database A
	dbA, err := sql.Open("godror", "userA/passwordA@dbA")
	if err != nil {
		log.Fatal(err)
	}
	defer dbA.Close()

	// Connection details for Oracle database B
	dbB, err := sql.Open("godror", "userB/passwordB@dbB")
	if err != nil {
		log.Fatal(err)
	}
	defer dbB.Close()

	// Amount to transfer
	amount := 100

	// Call the transferMoney function to perform the transaction
	if err := transferMoney(dbA, dbB, amount); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Money transferred successfully!")
}
