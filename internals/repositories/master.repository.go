// repositories/master_repo.go
package repositories

import (
	"context"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MasterRepo struct {
	DB *pgxpool.Pool
}

func NewMasterRepo(db *pgxpool.Pool) *MasterRepo {
	return &MasterRepo{DB: db}
}

func (r *MasterRepo) GetDirectors(ctx context.Context) ([]models.MasterDirector, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, name FROM directors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var directors []models.MasterDirector
	for rows.Next() {
		var d models.MasterDirector
		if err := rows.Scan(&d.ID, &d.Name); err != nil {
			return nil, err
		}
		directors = append(directors, d)
	}
	return directors, nil
}

func (r *MasterRepo) GetActors(ctx context.Context) ([]models.MasterActor, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, name FROM actors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var actors []models.MasterActor
	for rows.Next() {
		var a models.MasterActor
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}
		actors = append(actors, a)
	}
	return actors, nil
}

func (r *MasterRepo) GetLocations(ctx context.Context) ([]models.MasterLocation, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, name FROM location")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.MasterLocation
	for rows.Next() {
		var l models.MasterLocation
		if err := rows.Scan(&l.ID, &l.Name); err != nil {
			return nil, err
		}
		locations = append(locations, l)
	}
	return locations, nil
}

func (r *MasterRepo) GetTimes(ctx context.Context) ([]models.MasterTime, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, time FROM time")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var times []models.MasterTime
	for rows.Next() {
		var t models.MasterTime
		if err := rows.Scan(&t.ID, &t.Time); err != nil {
			return nil, err
		}
		times = append(times, t)
	}
	return times, nil
}

func (r *MasterRepo) GetCinemas(ctx context.Context) ([]models.MasterCinema, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, name, logo, price FROM cinema")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cinemas []models.MasterCinema
	for rows.Next() {
		var c models.MasterCinema
		if err := rows.Scan(&c.ID, &c.Name, &c.Logo, &c.Price); err != nil {
			return nil, err
		}
		cinemas = append(cinemas, c)
	}
	return cinemas, nil
}
