package database

import (
	"context"
	"math"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

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

func determineAlert(s *seedChild) bool {
	hasAlert := false

	if s.Saude != nil && len(s.Saude.Alerts) > 0 ||
		s.Educacao != nil && len(s.Educacao.Alerts) > 0 ||
		s.AssistenciaSocial != nil && len(s.AssistenciaSocial.Alerts) > 0 {
		hasAlert = true
	}

	return hasAlert
}

func insertAlert(ctx context.Context, pool *pgxpool.Pool, childID string, alert *alert) error {
	if alert == nil {
		return nil
	}

	var Alerts = map[string]string{
		"matricula_pendente":     "Matrícula Pendente",
		"frequencia_baixa":       "Frequência Baixa",
		"vacinas_atrasadas":      "Vacinas Atrasadas",
		"consulta_atrasada":      "Consulta Atrasada",
		"cadastro_ausente":       "Cadastro Ausente",
		"cadastro_desatualizado": "Cadastro Desatualizado",
		"beneficio_suspenso":     "Benefício Suspenso",
	}

	_, err := pool.Exec(ctx, `
		INSERT INTO alert (child_id, category, code, message, created_at)
		VALUES ($1, $2, $3, $4, NOW())`,
		childID, alert.Category, alert.Code, Alerts[alert.Code],
	)
	return err
}

func insertChild(ctx context.Context, pool *pgxpool.Pool, child *seedChild) error {
	age := computeAge(child.DataNascimento)
	hasAlert := determineAlert(child)

	var reviewedBy, reviewedAt *string
	if child.Revisado {
		reviewedBy = child.RevisadoPor
		reviewedAt = child.RevisadoEm
	}

	_, err := pool.Exec(ctx, `
		INSERT INTO children (id, name, age, neighborhood, has_alert, reviewed, reviewed_by, reviewed_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		ON CONFLICT (id) DO NOTHING`,
		child.ID, child.Nome, age, child.Bairro, hasAlert, child.Revisado, reviewedBy, reviewedAt,
	)
	return err
}

func insertHealth(ctx context.Context, pool *pgxpool.Pool, childID string, area *health) error {
	if area == nil {
		return nil
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO health (child_id, vaccinations_up_to_date, last_consultation, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (child_id) DO NOTHING`,
		childID, area.VaccinationsUpToDate, area.LastConsultation,
	)
	return err
}

func insertEducation(ctx context.Context, pool *pgxpool.Pool, childID string, area *education) error {
	if area == nil {
		return nil
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO education (child_id, school_name, frequency_percent, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (child_id) DO NOTHING`,
		childID, area.SchoolName, area.FrequenciaPercent,
	)
	return err
}

func insertSocialAssistance(ctx context.Context, pool *pgxpool.Pool, childID string, area *socialAssistance) error {
	if area == nil {
		return nil
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO social_assistance (child_id, cad_unico, active_benefit, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (child_id) DO NOTHING`,
		childID, area.CadUnico, area.BeneficioAtivo,
	)
	return err
}
