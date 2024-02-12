package sqlanalyzer

import (
	"strings"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gauth"
	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
	"github.com/Kwynto/GracefulDB/internal/engine/core"
	"github.com/Kwynto/GracefulDB/pkg/lib/e"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

// DCL — язык управления данными (Data Control Language)

func (q tQuery) DCLGrant() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLGrant"
	defer func() { e.Wrapper(op, err) }()

	return "DCLGrant", nil
}

func (q tQuery) DCLRevoke() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLRevoke"
	defer func() { e.Wrapper(op, err) }()

	return "DCLRevoke", nil
}

func (q tQuery) DCLUse() (result string, err error) {
	// -
	op := "internal -> analyzers -> sql -> DCL -> DCLUse"
	defer func() { e.Wrapper(op, err) }()

	// TODO: сделать проверку тикета и прав.

	db := core.RegExpCollection["UseWord"].ReplaceAllLiteralString(q.Instruction, " ")
	db = strings.TrimSpace(db)
	db = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(db, "")
	db = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(db, "")

	if !core.RegExpCollection["EntityName"].MatchString(db) {
		return ecowriter.EncodeString(gtypes.Response{
			State:  "error",
			Result: "invalid database name",
		}), nil
	}

	if core.LocalCoreSettings.FreezeMode {
		if _, ok := core.StorageInfo.DBs[db]; !ok {
			return ecowriter.EncodeString(gtypes.Response{
				State:  "error",
				Result: "the database does not exist",
			}), nil
		}
	}

	core.States[q.Ticket] = core.TState{
		CurrentDB: db,
	}

	return ecowriter.EncodeString(gtypes.Response{
		State:  "ok",
		Result: db,
	}), nil
}

func (q tQuery) DCLAuth() (result string, err error) {
	// This method is complete
	op := "internal -> analyzers -> sql -> DCL -> DCLAuth"
	defer func() { e.Wrapper(op, err) }()

	login := core.RegExpCollection["Login"].FindString(q.Instruction)
	login = core.RegExpCollection["LoginWord"].ReplaceAllLiteralString(login, " ")
	login = strings.TrimSpace(login)
	login = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(login, "")
	login = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(login, "")

	password := core.RegExpCollection["Password"].FindString(q.Instruction)
	password = core.RegExpCollection["PasswordWord"].ReplaceAllLiteralString(password, " ")
	password = strings.TrimSpace(password)
	password = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(password, "")
	password = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(password, "")

	hash := core.RegExpCollection["Hash"].FindString(q.Instruction)
	hash = core.RegExpCollection["HashWord"].ReplaceAllLiteralString(hash, " ")
	hash = strings.TrimSpace(hash)
	hash = core.RegExpCollection["QuotationMarks"].ReplaceAllLiteralString(hash, "")
	hash = core.RegExpCollection["SpecQuotationMark"].ReplaceAllLiteralString(hash, "")

	profile, err := gauth.GetProfile(login)
	if err != nil {
		return ecowriter.EncodeString(gtypes.Response{
			State: "auth error",
		}), nil
	}

	if profile.Status.IsBad() {
		return ecowriter.EncodeString(gtypes.Response{
			State: "auth error",
		}), nil
	}

	secret := gtypes.Secret{
		Login:    login,
		Password: password,
		Hash:     hash,
	}
	ticket, err := gauth.NewAuth(&secret)
	if err != nil {
		return ecowriter.EncodeString(gtypes.Response{
			State: "auth error",
		}), nil
	}

	return ecowriter.EncodeString(gtypes.Response{
		State:  "ok",
		Ticket: ticket,
	}), nil
}
