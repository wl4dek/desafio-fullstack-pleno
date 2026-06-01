"use client"

import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader } from "@/components/ui/card"
import { Skeleton } from "@/components/ui/skeleton"
import { useStatistic } from "@/hooks/useStatistic"
import { AlertCategoryType, AlertsCategories } from "@/types"
import { RefreshCw } from "lucide-react"
import { useState } from "react"
import { GeoJSON, MapContainer, TileLayer } from "react-leaflet"
import bairrosRio from '@/data/LimiteBairros.json'
import 'leaflet/dist/leaflet.css';

export function AlertsArea() {
    const { statistics, isLoading, isError, refresh } = useStatistic()
    const [metric, setMetric] = useState<AlertCategoryType>("health")
    const metrics = Object.keys(AlertsCategories);

    if (isLoading) {
        return (
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                {Array.from({ length: 4 }).map((_, i) => (
                    <Card key={i}>
                        <CardHeader className="pb-2">
                            <Skeleton className="h-4 w-24" />
                        </CardHeader>
                        <CardContent>
                            <Skeleton className="h-8 w-16" />
                        </CardContent>
                    </Card>
                ))}
            </div>
        )
    }

    if (isError) {
        return (
            <Card className="p-6 text-center">
                <p className="text-red-500 mb-2">Erro ao carregar indicadores</p>
                <Button variant="outline" size="sm" onClick={() => refresh()} className="gap-2">
                    <RefreshCw className="h-4 w-4" />
                    Tentar novamente
                </Button>
            </Card>
        )
    }

    if (!statistics || statistics.length === 0) return null

    const bairrosGeoJson = Object.fromEntries(
        statistics.map(item => [
            item.neighborhood.toLowerCase(),
            item
        ])
    );

    function getColor(value: number) {
        if (value >= 12) return "#800026";
        if (value >= 9) return "#BD0026";
        if (value >= 6) return "#E31A1C";
        if (value >= 3) return "#FC4E2A";
        if (value >= 1) return "#FD8D3C";
        if (value === 0) return "#e7e7e7";

        return "#FFEDA0";
    }

    return (
        <div>
            <div className="flex gap-2 mb-2">
                {metrics.map((key) => (
                    <button
                        key={key}
                        onClick={() => setMetric(key as AlertCategoryType)}
                        className={`px-3 py-1 rounded ${metric === key ? "bg-black text-white" : "bg-gray-200"
                            }`}
                    >
                        {AlertsCategories[key as AlertCategoryType]}
                    </button>
                ))}
            </div>
            <MapContainer className="h-120 w-full z-0"
                center={[-22.909996, -43.435081]}
                zoom={10}
            >
                <TileLayer
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />

                <GeoJSON
                    data={bairrosRio as any}
                    style={(feature) => {
                        const neighborhood = feature?.properties?.nome?.toLowerCase();
                        const info = bairrosGeoJson[neighborhood];
                        const value = info?.[metric] ?? 0;

                        return {
                            fillColor: getColor(value),
                            fillOpacity: 0.7,
                            color: "#444",
                            weight: 1,
                        };
                    }}
                    onEachFeature={(feature, layer) => {
                        const neighborhood = feature?.properties?.nome?.toLowerCase();
                        const info = bairrosGeoJson[neighborhood];
                        const value = info?.[metric] ?? 0;

                        layer.bindTooltip(
                            `<strong>${feature.properties.nome}</strong><br/>Valor: ${value}`,
                            { sticky: true }
                        );
                    }}
                />
            </MapContainer>
        </div>
    );
}