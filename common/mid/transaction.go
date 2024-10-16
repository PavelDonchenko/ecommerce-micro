package mid

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/PavelDonchenko/ecommerce-micro/common/database/sqldb"
	"github.com/PavelDonchenko/ecommerce-micro/common/httputils"
	"github.com/PavelDonchenko/ecommerce-micro/common/logger"
)

// BeginCommitRollback starts a transaction for the domain call.
func BeginCommitRollback(log *logger.Logger, bgn sqldb.Beginner) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		log.Info(ctx, "BEGIN TRANSACTION")
		tx, err := bgn.Begin()
		if err != nil {
			log.Error(ctx, "ROLLBACK TRANSACTION", "ERROR", err)
			httputils.NewInternalErrorWIthLog(c, err.Error(), log)
			return
		}

		defer func() {
			if err := recover(); err != nil {
				log.Info(ctx, "ROLLBACK TRANSACTION")
				encodedError, marshalError := json.Marshal(err)
				if marshalError != nil {
					encodedError = []byte("MARSHAL_ERROR")
				}
				if err = tx.Rollback(); err != nil {
					log.Info(ctx, "ROLLBACK TRANSACTION", "ERROR", string(encodedError))
				}
			}
		}()

		c.Request = c.Request.WithContext(setTran(ctx, tx))
		c.Next()

		if !StatusInList(c.Writer.Status(), []int{http.StatusOK, http.StatusCreated}) {
			log.Info(ctx, "ROLLBACK TRANSACTION")
			if err := tx.Rollback(); err != nil {
				if errors.Is(err, sql.ErrTxDone) {
					return
				}
				log.Info(ctx, "ROLLBACK TRANSACTION", "ERROR", err)
			}
		} else {
			log.Info(ctx, "COMMIT TRANSACTION")
			if err := tx.Commit(); err != nil {
				log.Error(ctx, "COMMIT TRANSACTION", "ERROR", err)
				httputils.NewInternalErrorWIthLog(c, err.Error(), log)
				return
			}
		}
	}
}

// StatusInList -> checks if the given status is in the list
func StatusInList(status int, statusList []int) bool {
	for _, i := range statusList {
		if i == status {
			return true
		}
	}
	return false
}
