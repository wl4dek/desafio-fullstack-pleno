package database

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"backend/internal/children"
)

type seedChild struct {
	ID                string            `json:"id"`
	Nome              string            `json:"nome"`
	DataNascimento    string            `json:"data_nascimento"`
	Bairro            string            `json:"bairro"`
	Responsavel       string            `json:"responsavel"`
	Saude             *health           `json:"saude"`
	Educacao          *education        `json:"educacao"`
	AssistenciaSocial *socialAssistance `json:"assistencia_social"`
	Revisado          bool              `json:"revisado"`
	RevisadoPor       *string           `json:"revisado_por"`
	RevisadoEm        *string           `json:"revisado_em"`
}

type socialAssistance struct {
	Alerts         []string `json:"alertas"`
	CadUnico       bool     `json:"cad_unico"`
	BeneficioAtivo bool     `json:"beneficio_ativo"`
}

type education struct {
	Alerts            []string `json:"alertas"`
	SchoolName        *string  `json:"escola"`
	FrequenciaPercent int      `json:"frequencia_percent"`
}

type health struct {
	Alerts               []string `json:"alertas"`
	LastConsultation     *string  `json:"ultima_consulta"`
	VaccinationsUpToDate bool     `json:"vacinas_em_dia"`
}

func LoadSeed(ctx context.Context, pool *pgxpool.Pool, path string) error {
	var count int
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM children").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check children count: %w", err)
	}

	if count > 0 {
		log.Info().Int("count", count).Msg("database already populated, skipping seed")
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read seed file: %w", err)
	}

	var seeds []seedChild
	if err := json.Unmarshal(data, &seeds); err != nil {
		return fmt.Errorf("failed to parse seed file: %w", err)
	}

	inserted := 0
	for _, s := range seeds {
		age := computeAge(s.DataNascimento)
		hasAlert := determineAlert(&s)

		var reviewedBy, reviewedAt *string
		if s.Revisado {
			reviewedBy = s.RevisadoPor
			reviewedAt = s.RevisadoEm
		}

		_, err := pool.Exec(ctx, `
			INSERT INTO children (id, name, age, neighborhood, has_alert, reviewed, reviewed_by, reviewed_at, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
			ON CONFLICT (id) DO NOTHING`,
			s.ID, s.Nome, age, s.Bairro, hasAlert, s.Revisado, reviewedBy, reviewedAt,
		)

		if err != nil {
			return fmt.Errorf("failed to insert seed child %s: %w", s.ID, err)
		}

		err = insertHealth(ctx, pool, s.ID, s.Saude)
		if err != nil {
			return fmt.Errorf("failed to insert health data for child %s: %w", s.ID, err)
		}

		err = insertEducation(ctx, pool, s.ID, s.Educacao)
		if err != nil {
			return fmt.Errorf("failed to insert education data for child %s: %w", s.ID, err)
		}

		err = insertSocialAssistance(ctx, pool, s.ID, s.AssistenciaSocial)
		if err != nil {
			return fmt.Errorf("failed to insert social assistance data for child %s: %w", s.ID, err)
		}
		inserted++
	}

	log.Info().Int("inserted", inserted).Msg("seed data loaded successfully")
	return nil
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

func determineAlert(s *seedChild) bool {
	hasAlert := false

	if s.Saude != nil && len(s.Saude.Alerts) > 0 ||
		s.Educacao != nil && len(s.Educacao.Alerts) > 0 ||
		s.AssistenciaSocial != nil && len(s.AssistenciaSocial.Alerts) > 0 {
		hasAlert = true
	}

	return hasAlert
}

func insertHealth(ctx context.Context, pool *pgxpool.Pool, childID string, area *health) error {
	if area == nil {
		return nil
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO health (child_id, vaccinations_up_to_date, alerts, last_consultation, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (child_id) DO NOTHING`,
		childID, area.VaccinationsUpToDate, area.Alerts, area.LastConsultation,
	)
	return err
}

func insertEducation(ctx context.Context, pool *pgxpool.Pool, childID string, area *education) error {
	if area == nil {
		return nil
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO education (child_id, school_name, frequency_percent, alerts, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (child_id) DO NOTHING`,
		childID, area.SchoolName, area.FrequenciaPercent, area.Alerts,
	)
	return err
}

func insertSocialAssistance(ctx context.Context, pool *pgxpool.Pool, childID string, area *socialAssistance) error {
	if area == nil {
		return nil
	}
	_, err := pool.Exec(ctx, `
		INSERT INTO social_assistance (child_id, cad_unico, active_benefit, alerts, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (child_id) DO NOTHING`,
		childID, area.CadUnico, area.BeneficioAtivo, area.Alerts,
	)
	return err
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}

// Ensure children package is used (for type compatibility in tests)
var _ = children.Child{}
