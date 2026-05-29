package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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

type alert struct {
	Code     string `json:"codigo"`
	Category string `json:"categoria"`
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
		err := insertChild(ctx, pool, &s)
		if err != nil {
			return fmt.Errorf("failed to insert seed child %s: %w", s.ID, err)
		}

		err = insertHealth(ctx, pool, s.ID, s.Saude)
		if err != nil {
			return fmt.Errorf("failed to insert health data for child %s: %w", s.ID, err)
		}
		if s.Saude != nil {
			for _, alertCode := range s.Saude.Alerts {
				err = insertAlert(ctx, pool, s.ID, &alert{
					Code:     alertCode,
					Category: "health",
				})
				if err != nil {
					return fmt.Errorf("failed to insert alert data for child %s: %w", s.ID, err)
				}
			}
		}

		err = insertEducation(ctx, pool, s.ID, s.Educacao)
		if err != nil {
			return fmt.Errorf("failed to insert education data for child %s: %w", s.ID, err)
		}
		if s.Educacao != nil {
			for _, alertCode := range s.Educacao.Alerts {
				err = insertAlert(ctx, pool, s.ID, &alert{
					Code:     alertCode,
					Category: "education",
				})
				if err != nil {
					return fmt.Errorf("failed to insert alert data for child %s: %w", s.ID, err)
				}
			}
		}

		err = insertSocialAssistance(ctx, pool, s.ID, s.AssistenciaSocial)
		if err != nil {
			return fmt.Errorf("failed to insert social assistance data for child %s: %w", s.ID, err)
		}
		if s.AssistenciaSocial != nil {
			for _, alertCode := range s.AssistenciaSocial.Alerts {
				err = insertAlert(ctx, pool, s.ID, &alert{
					Code:     alertCode,
					Category: "social_assistance",
				})
				if err != nil {
					return fmt.Errorf("failed to insert alert data for child %s: %w", s.ID, err)
				}
			}
		}
		inserted++
	}

	log.Info().Int("inserted", inserted).Msg("seed data loaded successfully")
	return nil
}

// Ensure children package is used (for type compatibility in tests)
var _ = children.Child{}
