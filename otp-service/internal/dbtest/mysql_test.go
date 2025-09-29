//go:build integration
// +build integration

package dbtest

import (
	"context"
	"database/sql"
	"runtime"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type MySQLContainerSuite struct {
	suite.Suite

	ctx    context.Context
	cancel context.CancelFunc
	mysql  *MySQLContainer
}

// -------- suite lifecycle --------

func (s *MySQLContainerSuite) SetupSuite() {
	// timeout รวมของชุดเทส
	s.ctx, s.cancel = context.WithTimeout(context.Background(), 3*time.Minute)

	// เริ่ม MySQL container
	mc, err := StartMySQL(
		s.ctx,
		WithMySQLDatabase("testdb"),
		WithMySQLUser("test"),
		WithMySQLPassword("test"),
		WithMySQLWaitDeadline(90*time.Second),
	)
	require.NoError(s.T(), err, "ควรสตาร์ทคอนเทนเนอร์ได้")
	s.mysql = mc
}

func (s *MySQLContainerSuite) TearDownSuite() {
	// stop container
	if s.mysql != nil {
		_ = s.mysql.Stop(context.Background())
	}
	// cancel context
	if s.cancel != nil {
		s.cancel()
	}
}

// -------- tests --------

func (s *MySQLContainerSuite) Test_ContainerFieldsAndDSN() {
	s.Require().NotNil(s.mysql, "container ต้องไม่เป็น nil")
	s.Require().NotEmpty(s.mysql.DSN, "DSN ต้องไม่ว่าง")
	s.Require().NotEmpty(s.mysql.Host, "Host ต้องไม่ว่าง")
	s.Require().NotEmpty(s.mysql.Port, "Port ต้องไม่ว่าง")
	s.Require().NotEmpty(s.mysql.User, "User ต้องไม่ว่าง")
	s.Require().NotEmpty(s.mysql.Password, "Password ต้องไม่ว่าง")
	s.Require().Equal("testdb", s.mysql.DB, "ชื่อ DB ควรเป็น testdb ตามที่ตั้ง")

	// parse DSN ด้วย mysql.ParseDSN (ของ go-sql-driver/mysql)
	cfg, err := mysql.ParseDSN(s.mysql.DSN)
	s.Require().NoError(err, "DSN ควร parse ได้")
	s.Equal("test", cfg.User, "user ควรเป็น test")
	s.Equal("test", cfg.Passwd, "password ควรเป็น test")
	s.Equal("testdb", cfg.DBName, "dbname ควรเป็น testdb")
	// ตัว driver จะเก็บ host:port ไว้ที่ cfg.Addr และ Net= "tcp"
	s.Equal("tcp", cfg.Net)
	s.Contains(cfg.Addr, ":", "addr ควรเป็น host:port")
}

func (s *MySQLContainerSuite) Test_PingAndSimpleQuery() {
	db, err := sql.Open("mysql", s.mysql.DSN)
	s.Require().NoError(err)
	defer db.Close()

	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)
	defer cancel()

	s.Require().NoError(db.PingContext(ctx), "ping DB ควรสำเร็จ")

	// DDL + DML + Query ง่าย ๆ
	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS t (
		id INT PRIMARY KEY AUTO_INCREMENT,
		v  INT NOT NULL
	)`)
	s.Require().NoError(err, "สร้างตารางควรสำเร็จ")

	_, err = db.ExecContext(ctx, `INSERT INTO t (v) VALUES (?)`, 42)
	s.Require().NoError(err, "insert ควรสำเร็จ")

	var got int
	err = db.QueryRowContext(ctx, `SELECT v FROM t LIMIT 1`).Scan(&got)
	s.Require().NoError(err, "query ควรสำเร็จ")
	s.Equal(42, got)
}

func (s *MySQLContainerSuite) Test_NewContainerWithCustomCreds() {
	// ทดสอบ start/stop อีกรอบด้วย credentials อื่น เพื่อยืนยันว่า options ทำงาน
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Minute)
	defer cancel()

	user, pass, dbname := "u2", "p2", "db2"
	mc2, err := StartMySQL(
		ctx,
		WithMySQLUser(user),
		WithMySQLPassword(pass),
		WithMySQLDatabase(dbname),
		WithMySQLWaitDeadline(60*time.Second),
	)
	s.Require().NoError(err)
	defer mc2.Stop(context.Background())

	cfg, err := mysql.ParseDSN(mc2.DSN)
	s.Require().NoError(err)
	s.Equal(user, cfg.User)
	s.Equal(pass, cfg.Passwd)
	s.Equal(dbname, cfg.DBName)
	s.Equal("tcp", cfg.Net)
}

func (s *MySQLContainerSuite) Test_Stop_ReleasesContainer() {
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Minute)
	defer cancel()

	tmp, err := StartMySQL(ctx, WithMySQLWaitDeadline(60*time.Second))
	s.Require().NoError(err)

	// ปิดคอนเทนเนอร์ควรไม่มี error
	err = tmp.Stop(context.Background())
	s.Require().NoError(err)

	// ไม่ assert ต่อว่าพอร์ตคืนทันที เพื่อเลี่ยง flakiness ระหว่าง OS/CI
}

// ตัวช่วยแสดง OS ที่รันอยู่ (มีประโยชน์เวลา debug CI)
func (s *MySQLContainerSuite) Test_RuntimeEnvInfo() {
	s.T().Logf("GOOS=%s GOARCH=%s", runtime.GOOS, runtime.GOARCH)
}

// -------- runner --------

func TestMySQLContainerSuite(t *testing.T) {
	suite.Run(t, new(MySQLContainerSuite))
}
