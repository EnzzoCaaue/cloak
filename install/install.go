package install

import (
	"fmt"
	"log"
	"sync"

	"github.com/raggaer/pigo"
)

var (
	cloakaTables = cloakaDatabase{
		pigo.Config.Key("database").String("database"),
		[]table{
			table{
				"cloaka_buypoints_paypal",
				`CREATE TABLE cloaka_buypoints_paypal (
                id int(11) NOT NULL AUTO_INCREMENT,
                payment_id varchar(255) DEFAULT NULL,
                state varchar(100) DEFAULT NULL,
                payer_email varchar(255) DEFAULT NULL,
                payer_first_name varchar(255) DEFAULT NULL,
                payer_second_name varchar(255) DEFAULT NULL,
                total int(11) DEFAULT NULL,
                points int(11) DEFAULT NULL,
                promo int(11) DEFAULT NULL,
                account int(11) DEFAULT NULL,
                PRIMARY KEY (id),
                UNIQUE KEY cloaka_buypoints_paypal_id_uindex (id)
                ) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;`,
			},
			table{
				"cloaka_news",
				`CREATE TABLE cloaka_news (
                id int(11) NOT NULL AUTO_INCREMENT,
                title varchar(45) DEFAULT NULL,
                text varchar(1400) DEFAULT NULL,
                created int(11) DEFAULT NULL,
                type int(11) DEFAULT NULL,
                PRIMARY KEY (id)
                ) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;`,
			},
			table{
				"cloaka_accounts",
				`CREATE TABLE cloaka_accounts (
                id int(11) NOT NULL AUTO_INCREMENT,
                account int(11) DEFAULT NULL,
                token varchar(40) DEFAULT '',
                points int(11) DEFAULT '0',
                admin int(11) DEFAULT '0',
                twofactor int(11) DEFAULT '0',
                recovery varchar(20) DEFAULT '',
                PRIMARY KEY (id),
                UNIQUE KEY cloaka_accounts_id_uindex (id),
                UNIQUE KEY cloaka_accounts_account_uindex (account)
                ) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;`,
			},
			table{
				"cloaka_online_records",
				`CREATE TABLE cloaka_online_records (
                id int(11) NOT NULL AUTO_INCREMENT,
                total int(11) DEFAULT NULL,
                at int(11) DEFAULT NULL,
                PRIMARY KEY (id),
                UNIQUE KEY cloaka_online_records_id_uindex (id)
                ) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;`,
			},
			table{
				"cloaka_players",
				`CREATE TABLE cloaka_players (
                id int(11) NOT NULL AUTO_INCREMENT,
                comment varchar(200) DEFAULT '',
                deleted int(11) DEFAULT '0',
                player_id int(11) DEFAULT NULL,
                signature varchar(100) DEFAULT '',
                hide int(11) DEFAULT '0',
                PRIMARY KEY (id),
                UNIQUE KEY cloaka_players_id_uindex (id),
                UNIQUE KEY cloaka_players_player_id_uindex (player_id)
                ) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;`,
			},
			table{
				"cloaka_shop_categories",
				`CREATE TABLE cloaka_shop_categories (
                id int(11) NOT NULL AUTO_INCREMENT,
                name varchar(50) DEFAULT NULL,
                description varchar(200) DEFAULT NULL,
                PRIMARY KEY (id),
                UNIQUE KEY cloaka_shop_categories_id_uindex (id)
                ) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;`,
			},
			table{
				"cloaka_shop_history",
				`CREATE TABLE cloaka_shop_history (
                id int(11) NOT NULL AUTO_INCREMENT,
                item int(11) DEFAULT NULL,
                status int(11) DEFAULT NULL,
                account int(11) DEFAULT NULL,
                PRIMARY KEY (id),
                UNIQUE KEY cloaka_shop_history_id_uindex (id)
                ) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;`,
			},
			table{
				"cloaka_shop_items",
				`CREATE TABLE cloaka_shop_items (
                id int(11) NOT NULL AUTO_INCREMENT,
                categorie_id int(11) DEFAULT NULL,
                item_id int(11) DEFAULT NULL,
                price int(11) DEFAULT NULL,
                name varchar(25) DEFAULT NULL,
                description varchar(200) DEFAULT NULL,
                PRIMARY KEY (id),
                UNIQUE KEY cloaka_shop_items_id_uindex (id)
                ) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=latin1;`,
			},
		},
	}
)

type cloakaDatabase struct {
	database string
	tables   []table
}

type table struct {
	name string
	sql  string
}

// Installer runs the installer to create the needed cloaka tables
func Installer() {
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(cloakaTables.tables))
	for _, t := range cloakaTables.tables {
		go func(t table) {
			if !t.isInstalled(cloakaTables.database) {
				fmt.Printf(" >> Installing missing table %v - ", t.name)
				if err := t.install(); err != nil {
					log.Fatal(err)
				}
				fmt.Printf("done \r\n")
			}
			waitGroup.Done()
		}(t)
	}
	waitGroup.Wait()
}

func (t *table) isInstalled(database string) bool {
	row := pigo.Database.QueryRow("SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_schema = ? AND table_name = ? LIMIT 1)", database, t.name)
	exists := false
	row.Scan(&exists)
	return exists
}

func (t *table) install() error {
	_, err := pigo.Database.Exec(t.sql)
	return err
}
