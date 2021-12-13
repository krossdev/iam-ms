// Copyright (c) 2021 Kross IAM Project.
// https://github.com/krossdev/iam/blob/main/LICENSE

package email

import (
	"github.com/krossdev/iam-ms/config"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// SMTP provider
type provider struct {
	Name    string `db:"name"`
	Host    string `db:"host"`
	Port    int    `db:"port"`
	SSLMode int    `db:"sslmode"`
	Sender  string `db:"sender"`
	Passwd  string `password`
}

func getSchemaPath(schema string) string {
	return filepath.Join(config.G.Mail.SchemaDir, schema)
}

func providersForSchema(schema string) (ps []*provider, err error) {
	var file []byte
	path := getSchemaPath(schema)
	if file, err = os.ReadFile(path); err != nil {
		err = errors.WithStack(err)
		return
	}
	var p *provider
	if err = yaml.Unmarshal(file, &p); err != nil {
		err = errors.WithStack(err)
		return
	}
	// Load smtp providers for specified schema
	ps = []*provider{p}
	return

	//q := fmt.Sprintf("select * from %s.smtp order by create_at", schema)
	//
	//DB := db.Default()
	//rows, err := DB.Queryx(DB.Rebind(q))
	//if err != nil {
	//	return nil, err
	//}
	//defer rows.Close()
	//
	//var providers []*provider
	//
	//for rows.Next() {
	//	provider := new(provider)
	//
	//	err := rows.StructScan(provider)
	//	if err != nil {
	//		return nil, err
	//	}
	//	providers = append(providers, provider)
	//	if len(providers) >= 20 {
	//		xlog.F("schema", schema).Warnf(
	//			"'%s' has too many smtp providers, only first 20 are used", schema,
	//		)
	//		break
	//	}
	//}
	//err = rows.Err()
	//if err != nil {
	//	return nil, err
	//}
	//if len(providers) == 0 {
	//	return nil, fmt.Errorf("'%s' has no smtp providers", schema)
	//}
	//return providers, nil
}
