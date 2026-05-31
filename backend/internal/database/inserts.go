package database

import (
	"context"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var alertMessages = map[string]string{
	"matricula_pendente":     "Matrícula Pendente",
	"frequencia_baixa":       "Frequência Baixa",
	"vacinas_atrasadas":      "Vacinas Atrasadas",
	"consulta_atrasada":      "Consulta Atrasada",
	"cadastro_ausente":       "Cadastro Ausente",
	"cadastro_desatualizado": "Cadastro Desatualizado",
	"beneficio_suspenso":     "Benefício Suspenso",
}

func computeAge(dateStr string) int {
	parsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return 0
	}
	years := int(math.Floor(time.Since(parsed).Hours() / 8760))
	if years < 0 {
		return 0
	}
	return years
}

func insertAlertHealth(ctx context.Context, pool *pgxpool.Pool, healthID int64, code string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO alert_health (health_id, code, message, created_at)
		VALUES ($1, $2, $3, NOW())`,
		healthID, code, alertMessages[code],
	)
	return err
}

func insertAlertEducation(ctx context.Context, pool *pgxpool.Pool, educationID int64, code string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO alert_education (education_id, code, message, created_at)
		VALUES ($1, $2, $3, NOW())`,
		educationID, code, alertMessages[code],
	)
	return err
}

func insertAlertSocialAssistance(ctx context.Context, pool *pgxpool.Pool, socialAssistanceID int64, code string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO alert_social_assistance (social_assistance_id, code, message, created_at)
		VALUES ($1, $2, $3, NOW())`,
		socialAssistanceID, code, alertMessages[code],
	)
	return err
}

func insertChild(ctx context.Context, pool *pgxpool.Pool, child *seedChild) error {
	age := computeAge(child.DataNascimento)

	var reviewedBy, reviewedAt *string
	if child.Revisado {
		reviewedBy = child.RevisadoPor
		reviewedAt = child.RevisadoEm
	}

	_, err := pool.Exec(ctx, `
		INSERT INTO children (id, name, age, neighborhood, reviewed, reviewed_by, reviewed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (id) DO NOTHING`,
		child.ID, child.Nome, age, child.Bairro, child.Revisado, reviewedBy, reviewedAt,
	)
	return err
}

func insertHealth(ctx context.Context, pool *pgxpool.Pool, childID string, area *health) (int64, error) {
	if area == nil {
		return 0, nil
	}
	var id int64
	err := pool.QueryRow(ctx, `
		INSERT INTO health (child_id, vaccinations_up_to_date, last_consultation, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (child_id) DO UPDATE SET child_id = EXCLUDED.child_id
		RETURNING id`,
		childID, area.VaccinationsUpToDate, area.LastConsultation,
	).Scan(&id)
	return id, err
}

func insertEducation(ctx context.Context, pool *pgxpool.Pool, childID string, area *education) (int64, error) {
	if area == nil {
		return 0, nil
	}
	var id int64
	err := pool.QueryRow(ctx, `
		INSERT INTO education (child_id, school_name, frequency_percent, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (child_id) DO UPDATE SET child_id = EXCLUDED.child_id
		RETURNING id`,
		childID, area.SchoolName, area.FrequenciaPercent,
	).Scan(&id)
	return id, err
}

func insertSocialAssistance(ctx context.Context, pool *pgxpool.Pool, childID string, area *socialAssistance) (int64, error) {
	if area == nil {
		return 0, nil
	}
	var id int64
	err := pool.QueryRow(ctx, `
		INSERT INTO social_assistance (child_id, cad_unico, active_benefit, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (child_id) DO UPDATE SET child_id = EXCLUDED.child_id
		RETURNING id`,
		childID, area.CadUnico, area.BeneficioAtivo,
	).Scan(&id)
	return id, err
}
