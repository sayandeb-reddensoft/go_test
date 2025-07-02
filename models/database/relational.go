package model

import "time"

// user
type Role struct {
	ID 		 uint   `gorm:"primaryKey"`
	RoleName string `gorm:"size:55,unique"`
}

type User struct {
	ID        uint   `gorm:"primaryKey"`
	UserId    string `gorm:"not null;unique"`
	Email     string `gorm:"size:255,unique"`
	Password  string `gorm:"size:255"`
	RoleID    uint  
	Role      Role   `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	IsActive  bool   `gorm:"default:false"`
}

type Login struct {
	Id              uint      `gorm:"primaryKey"`
	UserId          string    `gorm:"foreignKey:UserId,constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	LoginTimeStamp  time.Time `gorm:"default:null"`
	LogOutTimeStamp time.Time `gorm:"default:null"`
	LoginIp         string    `gorm:"size:50, default:null"`
	NoOfAttempts    int       `gorm:"default:null"`
}

// address
type Address struct {
	ID         		uint   `gorm:"primaryKey"`
	StreetLine 		string `gorm:"size:255"` 
	City       		string `gorm:"size:55"`
	State      		string `gorm:"size:55"`
	PostalIndexCode int32  `gorm:"default:null"`
}

// org
type Organization struct {
	ID        	  		uint    		`gorm:"primaryKey"`
	OrgName   	  		string  		`gorm:"size:255"`
	OrgUserID 	  		string		
	OrgUser   	  		User    		`gorm:"foreignKey:OrgUserID;references:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	AddressID 	  		uint		
	Address   	  		Address 		`gorm:"foreignKey:AddressID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	ContactName   		string  		`gorm:"size:55"`
	ContactNumber 		string  		`gorm:"size:55"`
	ContactDescription  string  		`gorm:"type:text"`
	CreatedAt 			time.Time  		`gorm:"autoCreateTime"`
	UpdatedAt 			time.Time  		`gorm:"autoUpdateTime"`
}
