package database

import "github.com/neo4j/neo4j-go-driver/neo4j"

func addContactTxFunc() neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		query := ""
		variables := map[string]interface{}{}
		return tx.Run(query, variables)
	}
}

func printContactsSummaryTxFunc() neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		var summary []string

		query := ""
		variables := map[string]interface{}{}
		result, err := tx.Run(query, variables)

		for result.Next() {
			summary = append(summary, result.Record().GetByIndex(0).(string))
		}

		if err = result.Err(); err != nil {
			return nil, err
		}

		return summary, err
	}
}

func printBondsSummaryTxFunc() neo4j.TransactionWork {
	return func(tx neo4j.Transaction) (interface{}, error) {
		var summary []string

		query := ""
		variables := map[string]interface{}{}
		result, err := tx.Run(query, variables)

		for result.Next() {
			summary = append(summary, result.Record().GetByIndex(0).(string))
		}

		if err = result.Err(); err != nil {
			return nil, err
		}

		return summary, err
	}
}
