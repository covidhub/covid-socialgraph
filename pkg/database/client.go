package database

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"go.uber.org/zap"
)

type Client struct {
	logger   *zap.SugaredLogger
	driver   neo4j.Driver
	url      string
	user     string
	password string
}

func New(logger *zap.SugaredLogger, url, user, password string) (*Client, error) {
	driver, err := neo4j.NewDriver(url, neo4j.BasicAuth(user, password, ""))
	if err != nil {
		return nil, err
	}

	client := &Client{logger: logger, driver: driver}
	return client, nil
}

func (c *Client) Close() error {
	return c.driver.Close()
}

func (c *Client) getProfile() error {
	return nil
}

func (c *Client) addProfile() error {
	return nil
}

func (c *Client) updateProfile() error {
	return nil
}

func (c *Client) getContacts() error {
	return nil
}

func (c *Client) getBonds() error {
	return nil
}

func (c *Client) addBond() error {
	return nil
}

func (c *Client) deleteBond() error {
	return nil
}

func (c *Client) addContact() error {
	var session neo4j.Session
	var err error

	if session, err = c.driver.Session(neo4j.AccessModeWrite); err != nil {
		return err
	}
	defer session.Close()

	_, err = session.WriteTransaction(addContactTxFunc())
	return err
}

func (c *Client) getReport() error {
	return nil
}

func (c *Client) getSummary() error {
	var session neo4j.Session
	var err error

	if session, err = c.driver.Session(neo4j.AccessModeRead); err != nil {
		return err
	}
	defer session.Close()

	contactSummary, err := session.ReadTransaction(printContactsSummaryTxFunc())
	if err != nil {
		return err
	}

	c.logger.Debugf("Result <getSummary>: %v", contactSummary)

	return err
}
