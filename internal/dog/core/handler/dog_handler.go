package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/wanrun-develop/wanrun/internal/dog/adapters/repository"
	"github.com/wanrun-develop/wanrun/internal/dog/core/dto"
	dwRepositoy "github.com/wanrun-develop/wanrun/internal/dogOwner/adapters/repository"
	model "github.com/wanrun-develop/wanrun/internal/models"
	"github.com/wanrun-develop/wanrun/pkg/errors"
	"github.com/wanrun-develop/wanrun/pkg/log"
	"github.com/wanrun-develop/wanrun/pkg/util"
)

type IDogHandler interface {
	GetAllDogs(echo.Context) ([]dto.DogListRes, error)
	GetDogByID(echo.Context, int64) (dto.DogDetailsRes, error)
	GetDogByDogOwnerID(echo.Context, int64) ([]dto.DogListRes, error)
	CreateDog(echo.Context, dto.DogSaveReq) (int64, error)
	UpdateDog(echo.Context, dto.DogSaveReq) (int64, error)
	DeleteDog(echo.Context, int64) error
}

type dogHandler struct {
	r   repository.IDogRepository
	dwr dwRepositoy.IDogOwnerRepository
}

func NewDogHandler(r repository.IDogRepository, dwr dwRepositoy.IDogOwnerRepository) IDogHandler {
	return &dogHandler{r, dwr}
}

func (h *dogHandler) GetAllDogs(c echo.Context) ([]dto.DogListRes, error) {
	logger := log.GetLogger(c).Sugar()

	dogs, err := h.r.GetAllDogs()

	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dog検索で失敗しました。", errors.NewDogServerErrorEType())
		return []dto.DogListRes{}, err
	}

	resDogs := []dto.DogListRes{}

	for _, d := range dogs {
		dr := dto.DogListRes{
			DogID:  d.DogID.Int64,
			Name:   d.Name.String,
			Weight: d.Weight.Int64,
			Sex:    d.Sex.String,
			Image:  d.Image.String,
			DogType: dto.DogTypeRes{
				DogTypeID: d.DogType.DogTypeID,
				Name:      d.DogType.Name,
			},
		}
		resDogs = append(resDogs, dr)
	}
	return resDogs, nil
}

// GetDogById: dogの詳細を検索して返す
//
// args:
//   - echo.Context:
//   - int64: 	dogのID
//
// return:
//   - dto.DogDetailsRes:	dogの詳細レスポンス
//   - error:	エラー
func (h *dogHandler) GetDogByID(c echo.Context, dogID int64) (dto.DogDetailsRes, error) {

	d, err := h.r.GetDogByID(dogID)

	if err != nil {
		return dto.DogDetailsRes{}, err
	}

	resDog := dto.DogDetailsRes{
		DogID:      d.DogID.Int64,
		DogOwnerID: d.DogOwnerID.Int64,
		Name:       d.Name.String,
		Weight:     d.Weight.Int64,
		Sex:        d.Sex.String,
		Image:      d.Image.String,
		CreateAt:   util.ConvertToWRTime(d.CreateAt),
		UpdateAt:   util.ConvertToWRTime(d.UpdateAt),
		DogType: dto.DogTypeRes{
			DogTypeID: d.DogType.DogTypeID,
			Name:      d.DogType.Name,
		},
	}
	return resDog, nil
}

// GetDogByDogOwnerID: dogの詳細を検索して返す
//
// args:
//   - echo.Context:	コンテキスト
//   - int64: 	dogOwnerのID
//
// return:
//   - []dto.DogListRes:	dogの一覧レスポンス
//   - error:	エラー
func (h *dogHandler) GetDogByDogOwnerID(c echo.Context, dogOwnerID int64) ([]dto.DogListRes, error) {
	logger := log.GetLogger(c).Sugar()

	logger.Infof("DogOwner %d の犬の一覧検索", dogOwnerID)

	//dogownerの検索（存在チェック)
	dogOwner, err := h.dwr.GetDogOwnerById(dogOwnerID)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogOwner検索で失敗しました。", errors.NewDogServerErrorEType())
		return []dto.DogListRes{}, err
	}
	if dogOwner.IsEmpty() {
		err = errors.NewWRError(nil, "指定されたdog ownerは存在しません。", errors.NewDogClientErrorEType())
		logger.Error("不正なdog owner idでの検索", err)
		return []dto.DogListRes{}, err
	}

	dogs, err := h.r.GetDogByDogOwnerID(dogOwner.DogOwnerID.Int64)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dog検索で失敗しました。", errors.NewDogServerErrorEType())
		return []dto.DogListRes{}, err
	}

	resDogs := []dto.DogListRes{}

	for _, d := range dogs {
		dr := dto.DogListRes{
			DogID:  d.DogID.Int64,
			Name:   d.Name.String,
			Weight: d.Weight.Int64,
			Sex:    d.Sex.String,
			Image:  d.Image.String,
			DogType: dto.DogTypeRes{
				DogTypeID: d.DogType.DogTypeID,
				Name:      d.DogType.Name,
			},
		}
		resDogs = append(resDogs, dr)
	}

	return resDogs, nil
}

// CreateDog: 犬の登録
//
//	dogownerの存在チェック
//
// args:
//   - echo.Context:	コンテキスト
//   - dto.DogSaveRew:	リクエスト内容
//
// return:
//   - int64:	登録したdogId
//   - error:	エラー
func (h *dogHandler) CreateDog(c echo.Context, saveReq dto.DogSaveReq) (int64, error) {
	logger := log.GetLogger(c).Sugar()

	logger.Info("create dog %v", saveReq)

	dogOwnerID := saveReq.DogOwnerID
	//dogownerの検索（存在チェック)
	if err := h.isExistsDogOwner(c, dogOwnerID); err != nil {
		return 0, err
	}

	dog := model.Dog{
		DogOwnerID: util.NewSqlNullInt64(dogOwnerID),
		Name:       util.NewSqlNullString(saveReq.Name),
		DogTypeID:  util.NewSqlNullInt64(saveReq.DogTypeID),
		Weight:     util.NewSqlNullInt64(saveReq.Weight),
		Sex:        util.NewSqlNullString(saveReq.Sex),
		Image:      util.NewSqlNullString(saveReq.Image),
	}

	dog, err := h.r.CreateDog(dog)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogOwnerの登録処理で失敗しました。", errors.NewDogServerErrorEType())
		return 0, err
	}

	return dog.DogID.Int64, err
}

// UpdateDog: dogの更新
//
//	dogの存在チェック
//
// args:
//   - echo.Context:	コンテキスト
//   - dto.DogSaveReq:	リクエスト内容
//
// return:
//   - dto.DogDetailsRes:	dog詳細用レスポンス
//   - error:	エラー
func (h *dogHandler) UpdateDog(c echo.Context, saveReq dto.DogSaveReq) (int64, error) {
	logger := log.GetLogger(c).Sugar()

	logger.Info("update dog %v", saveReq)

	dogID := saveReq.DogID
	// dogの存在チェック
	var dog model.Dog
	var err error
	if dog, err = h.isExistsDog(c, dogID); err != nil {
		return 0, err
	}

	//dogownerが変わっていれば存在チェック
	if saveReq.DogOwnerID != dog.DogOwnerID.Int64 {
		dogOwnerID := saveReq.DogOwnerID
		if err = h.isExistsDogOwner(c, dogOwnerID); err != nil {
			return 0, err
		}
	}

	//更新値をつめる
	dog.DogOwnerID = util.NewSqlNullInt64(saveReq.DogOwnerID)
	dog.Name = util.NewSqlNullString(saveReq.Name)
	dog.DogTypeID = util.NewSqlNullInt64(saveReq.DogTypeID)
	dog.Weight = util.NewSqlNullInt64(saveReq.Weight)
	dog.Sex = util.NewSqlNullString(saveReq.Sex)
	dog.Image = util.NewSqlNullString(saveReq.Image)
	//更新
	dog, err = h.r.UpdateDog(dog)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogOwnerの更新処理で失敗しました。", errors.NewDogServerErrorEType())
		return 0, err
	}

	return dog.DogID.Int64, err
}

func (h *dogHandler) DeleteDog(c echo.Context, dogID int64) error {
	if _, err := h.isExistsDog(c, dogID); err != nil {
		return err
	}
	if err := h.r.DeleteDog(dogID); err != nil {
		return err
	}
	return nil
}

// isExistsDog: dogの存在チェック
//
// args:
//   - echo.Context:	コンテキスト
//   - int64:	チェック対象のdogID
//
// return:
//   - error:	エラー
func (h *dogHandler) isExistsDog(c echo.Context, dogID int64) (model.Dog, error) {
	logger := log.GetLogger(c).Sugar()

	dog, err := h.r.GetDogByID(dogID)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dog検索で失敗しました。", errors.NewDogServerErrorEType())
		return model.Dog{}, err
	}
	if dog.IsEmpty() {
		err = errors.NewWRError(nil, "指定されたdogは存在しません。", errors.NewDogClientErrorEType())
		logger.Error("不正なdog owner idの指定", err)
		return model.Dog{}, err
	}
	return dog, nil
}

// isExistsDogOwner: dogOwnerの存在チェック
//
// args:
//   - echo.Context:	コンテキスト
//   - int64:	チェック対象のdogOwnerId
//
// return:
//   - error:	エラー
func (h *dogHandler) isExistsDogOwner(c echo.Context, dogOwnerID int64) error {
	logger := log.GetLogger(c).Sugar()
	//検索
	dogOwner, err := h.dwr.GetDogOwnerById(dogOwnerID)
	if err != nil {
		logger.Error(err)
		err = errors.NewWRError(err, "dogOwner検索で失敗しました。", errors.NewDogServerErrorEType())
		return err
	}
	if dogOwner.IsEmpty() {
		err = errors.NewWRError(nil, "指定されたdog ownerは存在しません。", errors.NewDogClientErrorEType())
		logger.Error("不正なdog owner idの指定", err)
		return err
	}
	return nil
}
