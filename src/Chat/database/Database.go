package database

import (
	"database/sql"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/common"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/Event"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/Password"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/Session"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/User"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/entities/persistent/room"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/enum"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/models/request"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/Chat/settings"
	"github.com/vault-thirteen/Simpel-Chat-Server/src/helper"
)

type Database struct {
	criticalErrorsChan *chan error
	sqlDb              *sql.DB
	gormDb             *gorm.DB
}

func NewDatabase(settings *settings.ChatDatabaseSettings, criticalErrorsChan *chan error) (db *Database, err error) {
	switch settings.Type {
	case enum.DatabaseType_MySQL:
		return NewMysqlDatabase(settings, criticalErrorsChan)
	default:
		return nil, errors.New("unsupported database type: " + settings.Type.ToString())
	}
}
func NewMysqlDatabase(settings *settings.ChatDatabaseSettings, criticalErrorsChan *chan error) (db *Database, err error) {
	mc := mysql.Config{
		Net:                  settings.NetworkType,
		Addr:                 net.JoinHostPort(settings.HostName, strconv.Itoa(int(settings.PortNumber))),
		DBName:               settings.DBName,
		User:                 settings.DBUserName,
		Passwd:               settings.DBPassword,
		AllowNativePasswords: settings.AllowNativePasswords,
		CheckConnLiveness:    settings.CheckConnLiveness,
		MaxAllowedPacket:     settings.MaxAllowedPacket,
		Params:               settings.Parameters,
	}

	db = new(Database)

	db.criticalErrorsChan = criticalErrorsChan

	db.sqlDb, err = sql.Open(settings.DriverName, mc.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = db.sqlDb.Ping()
	if err != nil {
		return nil, err
	}

	db.gormDb, err = gorm.Open(gmysql.New(gmysql.Config{Conn: db.sqlDb}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if settings.UseDataInitialisation {
		classesToInit := []any{
			&rq.ChangePassword{}, // Request for password change.
			&ev.Event{},          // Event.
			&rq.LogIn{},          // Request for logging in.
			&rq.LogOut{},         // Request for logging out.
			&pwd.Password{},      // User's password.
			&rq.Registration{},   // Request for registration.
			&rm.Room{},           // Chat room.
			&ses.Session{},       // User's session.
			&usr.User{},          // User.
		}

		for _, cti := range classesToInit {
			err = db.gormDb.AutoMigrate(cti)
			if err != nil {
				return nil, err
			}
		}
	}

	return db, nil
}
func (db *Database) Close() (err error) {
	var sqlDb *sql.DB
	sqlDb, err = db.gormDb.DB()
	if err != nil {
		return err
	}

	err = sqlDb.Close()
	if err != nil {
		return err
	}

	err = db.sqlDb.Close()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) CleanAllSessions() (err error) {
	// Remove all records keeping the auto-increment counters (ID, ...).
	tx := db.gormDb.Session(&gorm.Session{AllowGlobalUpdate: true}).Unscoped().Delete(&ses.Session{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (db *Database) ListAllRooms() (rooms []*rm.Room, err error) {
	tx := db.gormDb.Model(&rm.Room{}).Find(&rooms)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return rooms, nil
}
func (db *Database) ListAdministratorUsers(ids common.IdList) (administratorUsers []*usr.User, err error) {
	tx := db.gormDb.Model(&usr.User{}).Where("id IN ?", ids.AsArray()).Find(&administratorUsers)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return administratorUsers, nil
}
func (db *Database) IsEmailAddressUsed(email string) (isUsed bool, err error) {
	var numOfRegistrations int64
	tx := db.gormDb.Model(&rq.Registration{}).Where("email = ?", email).Count(&numOfRegistrations)
	if tx.Error != nil {
		return false, tx.Error
	}

	var numOfUsers int64
	tx = db.gormDb.Model(&usr.User{}).Where("email = ?", email).Count(&numOfUsers)
	if tx.Error != nil {
		return false, tx.Error
	}

	return (numOfRegistrations + numOfUsers) > 0, nil
}
func (db *Database) CreateRegistrationRequest(rr *rq.Registration) (err error) {
	tx := db.gormDb.Create(&rr)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) GetFirstOutdatedRegistrationRequest(edgeTime time.Time) (rrs []rq.Registration, err error) {
	tx := db.gormDb.Limit(1).Where("created_at <= ?", edgeTime).Find(&rrs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return rrs, nil
}
func (db *Database) DeleteRegistrationRequest(rr *rq.Registration) (err error) {
	tx := db.gormDb.Delete(&rr)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) FindRegistrationRequest(rr *rq.Registration) (err error) {
	tx := db.gormDb.First(rr, "requestId = ?", rr.RequestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (db *Database) GetRegistrationRequestByEmail(userEmail string, rr *rq.Registration) (err error) {
	tx := db.gormDb.First(rr, "email = ?", userEmail)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (db *Database) CreateUser(user *usr.User, password string) (err error) {
	var p = &pwd.Password{
		UserId: user.Id,
		Text:   password,
	}

	user.Password = p

	tx := db.gormDb.Create(user)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) ExistsUserWithEmail(email string) (exists bool, err error) {
	var numOfUsers int64
	tx := db.gormDb.Model(&usr.User{}).Where("email = ?", email).Count(&numOfUsers)
	if tx.Error != nil {
		return false, tx.Error
	}

	return numOfUsers > 0, nil
}
func (db *Database) ExistsUserWithId(user *usr.User) (exists bool, err error) {
	var numOfUsers int64
	tx := db.gormDb.Model(&usr.User{}).Where("id = ?", user.Id).Count(&numOfUsers)
	if tx.Error != nil {
		return false, tx.Error
	}

	return numOfUsers > 0, nil
}
func (db *Database) IsUserWithEmailLoggedIn(email string) (isLoggedIn bool, err error) {
	var numOfUsers int64
	tx := db.gormDb.Model(&ses.Session{}).Joins("left join users on users.id = sessions.user_id").
		Where("users.email = ?", email).Count(&numOfUsers)
	if tx.Error != nil {
		return false, tx.Error
	}

	return numOfUsers > 0, nil
}
func (db *Database) IsUserWithEmailBanned(email string) (isBanned bool, err error) {
	var numOfUsers int64
	tx := db.gormDb.Model(&usr.User{}).Where("email = ?", email).Where("isBanned = ?", true).Count(&numOfUsers)
	if tx.Error != nil {
		return false, tx.Error
	}

	return numOfUsers > 0, nil
}
func (db *Database) CreateLogInRequest(lir *rq.LogIn) (err error) {
	tx := db.gormDb.Create(&lir)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) GetFirstOutdatedLogInRequest(edgeTime time.Time) (lirs []rq.LogIn, err error) {
	tx := db.gormDb.Limit(1).Where("created_at <= ?", edgeTime).Find(&lirs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return lirs, nil
}
func (db *Database) DeleteLogInRequest(lir *rq.LogIn) (err error) {
	tx := db.gormDb.Delete(&lir)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) FindLogInRequest(lir *rq.LogIn) (err error) {
	tx := db.gormDb.First(lir, "requestId = ?", lir.RequestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (db *Database) CheckUserPassword(user *usr.User, pwd *string) (ok bool, err error) {
	var numOfUsers int64
	tx := db.gormDb.Model(&usr.User{}).Joins("left join passwords on passwords.user_id = users.id").
		Where("users.email = ?", user.EmailAddress).Where("passwords.text = ?", *pwd).Count(&numOfUsers)
	if tx.Error != nil {
		return false, tx.Error
	}

	return numOfUsers == 1, nil
}
func (db *Database) GetLogInRequestByEmail(userEmail string, lir *rq.LogIn) (err error) {
	tx := db.gormDb.First(lir, "email = ?", userEmail)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (db *Database) CreateSession(session *ses.Session) (err error) {
	tx := db.gormDb.Create(session)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) CountAllSessions() (n int, err error) {
	var count int64
	tx := db.gormDb.Model(&ses.Session{}).Count(&count)
	if tx.Error != nil {
		return -1, tx.Error
	}

	return int(count), nil
}
func (db *Database) FindUserWithEmail(user *usr.User) (err error) {
	tx := db.gormDb.First(user, "email = ?", user.EmailAddress)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (db *Database) CreateLogOutRequest(lor *rq.LogOut) (err error) {
	tx := db.gormDb.Create(&lor)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) GetLogOutRequestByRequestId(requestId string, lor *rq.LogOut) (err error) {
	tx := db.gormDb.First(lor, "requestId = ?", requestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (db *Database) DeleteLogOutRequest(lor *rq.LogOut) (err error) {
	tx := db.gormDb.Delete(&lor)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) DeleteSession(session *ses.Session) (err error) {
	tx := db.gormDb.Delete(&session)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) GetUserById(user *usr.User) (err error) {
	tx := db.gormDb.First(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (db *Database) GetUserWithPasswordById(user *usr.User) (err error) {
	tx := db.gormDb.Joins("Password").First(&user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (db *Database) CreatePasswordChangeRequest(pcr *rq.ChangePassword) (err error) {
	tx := db.gormDb.Create(&pcr)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) FindPasswordChangeRequest(pcr *rq.ChangePassword) (err error) {
	tx := db.gormDb.First(pcr, "requestId = ?", pcr.RequestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (db *Database) DeletePasswordChangeRequest(pcr *rq.ChangePassword) (err error) {
	tx := db.gormDb.Delete(&pcr)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) ChangeUserPassword(user *usr.User, newPassword *string) (err error) {
	user.Password.Text = *newPassword

	tx := db.gormDb.Session(&gorm.Session{FullSaveAssociations: true}).Save(&user)
	if tx.Error != nil {
		return tx.Error
	}

	// New password must be different,
	// so we should get a normal value for RowsAffected.
	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) GetFirstOutdatedLogOutRequest(edgeTime time.Time) (lors []rq.LogOut, err error) {
	tx := db.gormDb.Limit(1).Where("created_at <= ?", edgeTime).Find(&lors)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return lors, nil
}
func (db *Database) GetFirstOutdatedPasswordChangeRequest(edgeTime time.Time) (pcrs []rq.ChangePassword, err error) {
	tx := db.gormDb.Limit(1).Where("created_at <= ?", edgeTime).Find(&pcrs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return pcrs, nil
}
func (db *Database) BanUserById(user *usr.User) (err error) {
	tx := db.gormDb.Model(&user).Update("isBanned", true)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) CreateEvent(event *ev.Event) (err error) {
	if event == nil {
		return errors.New(helper.Err_NullPointer)
	}

	if !event.HasValidType() {
		return helper.NewError_InvalidEnumValue(enum.EnumField_EventType, event.Type)
	}

	tx := db.gormDb.Create(event)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) CreateRoom(room *rm.Room) (err error) {
	if room == nil {
		return errors.New(helper.Err_NullPointer)
	}

	err = room.Validate()
	if err != nil {
		return err
	}

	tx := db.gormDb.Create(room)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) CountAllRooms() (n int, err error) {
	var count int64
	tx := db.gormDb.Model(&rm.Room{}).Count(&count)
	if tx.Error != nil {
		return -1, tx.Error
	}

	return int(count), nil
}
func (db *Database) DeleteRoom(room *rm.Room) (err error) {
	if room == nil {
		return errors.New(helper.Err_NullPointer)
	}

	if room.Id == 0 {
		return errors.New(helper.Err_IdIsNotSet)
	}

	tx := db.gormDb.Delete(&room)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) SaveRoom(room *rm.Room) (err error) {
	if room == nil {
		return errors.New(helper.Err_NullPointer)
	}

	if room.Id == 0 {
		return errors.New(helper.Err_IdIsNotSet)
	}

	tx := db.gormDb.Save(room)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New(helper.Err_NoRowsAreAffected)
	}

	return nil
}
func (db *Database) CountAllUsers() (n int, err error) {
	var count int64
	tx := db.gormDb.Model(&usr.User{}).Count(&count)
	if tx.Error != nil {
		return -1, tx.Error
	}

	return int(count), nil
}
func (db *Database) ListUsers(pageSize int, pageNumber int) (users []*usr.User, err error) {
	tx := db.gormDb.Order("id asc").Limit(pageSize).Offset((pageNumber - 1) * pageSize).Find(&users)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return users, nil
}
