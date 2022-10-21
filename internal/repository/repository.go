package repository

type DatabaseRepo interface {
	AddUser(res models.user) err
}