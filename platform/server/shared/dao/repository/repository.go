/*
 *
 *  * Copyright 2022 CloudWeGo Authors
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  *
 *  *     http://www.apache.org/licenses/LICENSE-2.0
 *  *
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 */

package repository

import (
	"errors"
	"github.com/cloudwego/cwgo/platform/server/shared/consts"
	"github.com/cloudwego/cwgo/platform/server/shared/kitex_gen/repository"
	"github.com/cloudwego/cwgo/platform/server/shared/utils"
	"gorm.io/gorm"
)

type IRepositoryDaoManager interface {
	GetTokenByID(id int64) (string, error)
	GetRepoTypeByID(id int64) (int32, error)
	GetRepository(id int64) (*repository.Repository, error)
	ChangeRepositoryStatus(id int64, status string) error
	GetAllRepositories() ([]*repository.Repository, error)

	AddRepository(repoURL, token, status string, repoType int32) error
	DeleteRepository(ids []string) error
	UpdateRepository(id, token, status string) error
	GetRepositories(page, limit, order int32, orderBy string) ([]*repository.Repository, error)
}

type MysqlRepositoryManager struct {
	db *gorm.DB
}

var _ IRepositoryDaoManager = (*MysqlRepositoryManager)(nil)

func NewMysqlRepository(db *gorm.DB) *MysqlRepositoryManager {
	return &MysqlRepositoryManager{
		db: db,
	}
}

func (r *MysqlRepositoryManager) GetTokenByID(id int64) (string, error) {
	var repo repository.Repository
	result := r.db.Table(consts.TableNameRepository).Model(&repo).Where("id = ?", id).Take(&repo)
	if result.Error != nil {
		return "", result.Error
	}

	return repo.Token, nil
}

func (r *MysqlRepositoryManager) GetRepoTypeByID(id int64) (int32, error) {
	var repo repository.Repository
	result := r.db.Table(consts.TableNameRepository).Model(&repo).Where("id = ?", id).Take(&repo)
	if result.Error != nil {
		return 0, result.Error
	}

	return repo.RepoType, nil
}

func (r *MysqlRepositoryManager) GetRepository(id int64) (*repository.Repository, error) {
	var repo *repository.Repository

	result := r.db.Table(consts.TableNameRepository).Where("id = ?", id).Find(&repo)
	if result.Error != nil {
		return nil, result.Error
	}

	return repo, nil
}

func (r *MysqlRepositoryManager) ChangeRepositoryStatus(id int64, status string) error {
	if !utils.ValidStatus(status) {
		return errors.New("invalid status")
	}
	timeNow := utils.GetCurrentTime()
	result := r.db.Table(consts.TableNameRepository).Model(&repository.Repository{}).Where("id = ?", id).Updates(repository.Repository{
		Status:     status,
		UpdateTime: timeNow,
	})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepositoryManager) GetAllRepositories() ([]*repository.Repository, error) {
	var repos []*repository.Repository
	result := r.db.Table(consts.TableNameRepository).Find(&repos)
	if result.Error != nil {
		return nil, result.Error
	}

	return repos, nil
}

func (r *MysqlRepositoryManager) AddRepository(repoURL, token, status string, repoType int32) error {
	timeNow := utils.GetCurrentTime()
	repo := repository.Repository{
		RepositoryUrl:  repoURL,
		Token:          token,
		Status:         status,
		RepoType:       repoType,
		LastUpdateTime: "0",
		LastSyncTime:   timeNow,
		CreateTime:     timeNow,
		UpdateTime:     timeNow,
	}
	result := r.db.Table(consts.TableNameRepository).Create(&repo)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepositoryManager) DeleteRepository(ids []string) error {
	var repo repository.Repository
	result := r.db.Table(consts.TableNameRepository).Delete(&repo, ids)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepositoryManager) UpdateRepository(id, token, status string) error {
	if !utils.ValidStatus(status) {
		return errors.New("invalid status")
	}
	timeNow := utils.GetCurrentTime()
	result := r.db.Table(consts.TableNameRepository).Model(&repository.Repository{}).Where("id = ?", id).Updates(repository.Repository{
		Token:      token,
		UpdateTime: timeNow,
		Status:     status,
	})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *MysqlRepositoryManager) GetRepositories(page, limit, order int32, orderBy string) ([]*repository.Repository, error) {
	var repos []*repository.Repository

	// Default sort field to 'update_time' if not provided
	if orderBy == "" {
		orderBy = consts.OrderByUpdateTime
	}

	switch order {
	case consts.OrderNumInc:
		orderBy = orderBy + " " + consts.Inc
	case consts.OrderNumDec:
		orderBy = orderBy + " " + consts.Dec
	default:
		orderBy = orderBy + " " + consts.Inc
	}

	offset := (page - 1) * limit
	result := r.db.Table(consts.TableNameRepository).Offset(int(offset)).Limit(int(limit)).Order(orderBy).Find(&repos)
	if result.Error != nil {
		return nil, result.Error
	}

	return repos, nil
}
