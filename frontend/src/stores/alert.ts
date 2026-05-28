export const Alerts = {
    'matricula_pendente': "Matrícula Pendente",
    'frequencia_baixa': "Frequência Baixa",
    'vacinas_atrasadas': "Vacinas Atrasadas",
    'consulta_atrasada': "Consulta Atrasada",
    'cadastro_ausente': "Cadastro Ausente",
    'cadastro_desatualizado': "Cadastro Desatualizado",
    'beneficio_suspenso': "Benefício Suspenso",
} as const

export type AlertType = keyof typeof Alerts

