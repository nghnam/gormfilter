# gormfilter - dynamic filter like django-filter

Work with Gorm v2.

# Installation
```sh
go get github.com/nghnam/gormfilter
```


# Example:

```go
package main

import (
        "time"

        "gorm.io/driver/sqlite"
        "gorm.io/gorm"
        "gorm.io/gorm/logger"

        "github.com/nghnam/gormfilter"
)

type Role struct {
        ID   uint   `gorm:"primarykey"`
        Name string `gorm:"type:varchar(15)"`
}

func (Role) TableName() string {
        return "role"
}

type User struct {
        gorm.Model
        Role     Role   `gorm:"foreignkey:RoleID" json:"-" binding:"omitempty"`
        RoleID   uint   `json:"role_id"`
        FullName string `gorm:"type:varchar(50)" json:"full_name"`
        Phone    string `gorm:"type:varchar(20)" json:"phone"`
}

func (User) TableName() string {
        return "user"
}

type UserFilter struct {
        FullName        *string    `filter:"column:full_name;op:contains"`
        Phone           *string    `filter:"column:phone;op:contains"`
        CreatedAtAfter  *time.Time `filter:"column:created_at;op:gte" time_format:"2006-01-02 15:04:05"`
        CreatedAtBefore *time.Time `filter:"column:created_at;op:lte" time_format:"2006-01-02 15:04:05"`
        RoleName        []string   `filter:"column:Role.name;op:!in"` // Join
}

func main() {
        db, _ := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{
                Logger: logger.Default.LogMode(logger.Info),
        })
        db.AutoMigrate(Role{})
        db.AutoMigrate(User{})

        timeForm := "2006-01-02 15:04:05"
        fullName := "Nam"
        phone := "555"
        createdAtAfter, _ := time.Parse(timeForm, "2020-01-10 00:00:00")
        createdAtBefore, _ := time.Parse(timeForm, "2020-01-20 00:00:00")
        uf := &UserFilter{
                FullName:        &fullName,
                Phone:           &phone,
                CreatedAtAfter:  &createdAtAfter,
                CreatedAtBefore: &createdAtBefore,
                RoleName:        []string{"Manager", "Staff"},
        }

        var users []User
        db, _ = gormfilter.FilterQuery(db.Model(&User{}), uf)
        _ = db.Find(&users).Error
}
```

```
[0.189ms] [rows:0] SELECT `user`.`id`,`user`.`created_at`,`user`.`updated_at`,`user`.`deleted_at`,`user`.`role_id`,`user`.`full_name`,`user`.`phone`,`Role`.`id` AS `Role__id`,`Role`.`name` AS `Role__name` FROM `user` LEFT JOIN `role` `Role` ON `user`.`role_id` = `Role`.`id` WHERE `full_name` LIKE "%Nam%" AND `phone` LIKE "%555%" AND `created_at` >= "2020-01-10 00:00:00" AND `created_at` <= "2020-01-20 00:00:00" AND `Role`.`name` NOT IN ("Manager","Staff") AND `user`.`deleted_at` IS NULL
```

# Expressions
See `expression.go`

# Note
- sqlite does not support `ILIKE` expression