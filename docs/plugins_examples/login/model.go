//////////////////////////////////////////////////////////////////////////////////
// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>                          //
//                                                                              //
// This file is part of the ERP library.                                        //
//                                                                              //
//  Unauthorized copying of this file, via any media is strictly prohibited     //
//  Proprietary and confidential                                                //
//////////////////////////////////////////////////////////////////////////////////

package model

import (
	"time"

	"gorm.io/gorm"
)

type UsersInfo struct {
	gorm.Model
	Update      time.Time `gorm:"column:update;type:timestamp;DEFAULT CURRENT_TIMESTAMP;"`
	UID         string    `gorm:"column:uid;type:varchar(60);UNIQUE;NOT NULL;"` //userID
	Name        string    `gorm:"column:name;type:varchar(254);DEFAULT '';"`
	MName       string    `gorm:"column:mname;type:varchar(254);DEFAULT '';"`
	Surname     string    `gorm:"column:surname;type:varchar(254);DEFAULT '';"`
	AvatarID    string    `gorm:"column:avatarid;type:varchar(60);DEFAULT '';"` //fileid in files table
	BirthDate   time.Time `gorm:"column:birthdate;type:timestamp;DEFAULT '0';"`
	PhoneNumber string    `gorm:"column:phonenumber;type:varchar(254);DEFAULT '';"`
	CompanyID   string    `gorm:"column:companyid;type:varchar(254);DEFAULT '';"`
	StaffID     string    `gorm:"column:staffid;type:varchar(254);DEFAULT '';"`
	TFA         int       `gorm:"column:tfa;type:double;DEFAULT false;"`
	TFAType     string    `gorm:"column:tfatype;type:varchar(60);DEFAULT '';"`
	DateFormat  string    `gorm:"column:dateformat;type:varchar(60);DEFAULT '2006-01-02';"`
	Companies   int       `gorm:"column:tfa;type:int;DEFAULT 0;"` //How many companies can user create
}
