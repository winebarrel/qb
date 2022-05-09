package qb

import (
	"fmt"
)

var stmtBldrByName = map[string]func(scale int) []string{
	"tpcb-like": func(scale int) []string {
		return []string{
			fmt.Sprintf("SET @aid = CEIL(RAND() * %d)", NAccounts*scale),
			fmt.Sprintf("SET @bid = CEIL(RAND() * %d)", NBranches*scale),
			fmt.Sprintf("SET @tid = CEIL(RAND() * %d)", NTellers*scale),
			"SET @delta = FLOOR(RAND() * 10001) - 5000",
			"BEGIN",
			"UPDATE qb_accounts SET abalance = abalance + @delta WHERE aid = @aid",
			"SELECT abalance FROM qb_accounts WHERE aid = @aid",
			"UPDATE qb_tellers SET tbalance = tbalance + @delta WHERE tid = @tid",
			"UPDATE qb_branches SET bbalance = bbalance + @delta WHERE bid = @bid",
			"INSERT INTO qb_history (tid, bid, aid, delta, mtime) VALUES (@tid, @bid, @aid, @delta, CURRENT_TIMESTAMP)",
			"COMMIT",
		}
	},
	"simple-update": func(scale int) []string {
		return []string{
			fmt.Sprintf("SET @aid = CEIL(RAND() * %d)", NAccounts*scale),
			fmt.Sprintf("SET @bid = CEIL(RAND() * %d)", NBranches*scale),
			fmt.Sprintf("SET @tid = CEIL(RAND() * %d)", NTellers*scale),
			"SET @delta = FLOOR(RAND() * 10001) - 5000",
			"BEGIN",
			"UPDATE qb_accounts SET abalance = abalance + @delta WHERE aid = @aid",
			"SELECT abalance FROM qb_accounts WHERE aid = @aid",
			"INSERT INTO qb_history (tid, bid, aid, delta, mtime) VALUES (@tid, @bid, @aid, @delta, CURRENT_TIMESTAMP)",
			"COMMIT",
		}
	},
	"select-only": func(scale int) []string {
		return []string{
			fmt.Sprintf("SET @aid = CEIL(RAND() * %d)", NAccounts*scale),
			"SELECT abalance FROM qb_accounts WHERE aid = @aid",
		}
	},
}

func ScriptNames() []string {
	names := []string{}

	for name := range stmtBldrByName {
		names = append(names, name)
	}

	return names
}

func NewScript(name string, scale int) ([]string, error) {
	bldr, ok := stmtBldrByName[name]

	if !ok {
		return nil, fmt.Errorf("script not found: %s", name)
	}

	return bldr(scale), nil
}
