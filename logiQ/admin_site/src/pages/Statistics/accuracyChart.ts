import type {EChartsOption} from 'echarts';

import type {AccuracyByDimension, AccuracyDistributions} from './services';
import {DIMENSION_KEYS, DIMENSION_LABELS, getDimensionValueLabel} from './dimensionLabels';
import {getDimensionColor} from './chartColors';

const MAX_BARS_PER_GROUP = 3;

export interface AccuracyChartItem {
    label: string;
    accuracy: number;
}

export interface AccuracyChartGroup {
    dimensionKey: (typeof DIMENSION_KEYS)[number];
    dimensionLabel: string;
    items: AccuracyChartItem[];
}

export interface AccuracyChartPoint {
    value: number;
    accuracy: number;
    dimensionLabel: string;
    valueLabel: string;
    itemStyle?: {
        color?: string;
        borderRadius?: [number, number, number, number];
    };
}

export const toAccuracyChartGroups = (accuracyDistributions?: AccuracyDistributions): AccuracyChartGroup[] => {
    if (!accuracyDistributions) {
        return [];
    }

    return DIMENSION_KEYS.map((key) => {
        const items: AccuracyByDimension[] = accuracyDistributions[key] ?? [];
        const groupItems = items.slice(0, MAX_BARS_PER_GROUP).map(({value, accuracy}) => ({
            label: getDimensionValueLabel(key, value),
            accuracy,
        }));

        return {
            dimensionKey: key,
            dimensionLabel: DIMENSION_LABELS[key],
            items: groupItems,
        } satisfies AccuracyChartGroup;
    }).filter((group) => group.items.length > 0);
};

export const buildAccuracyOption = (groups: AccuracyChartGroup[]): EChartsOption => {
    const xAxisCategories = groups.map((group) => group.dimensionLabel);

    const maxItems = groups.reduce((max, group) => Math.max(max, group.items.length), 0);
    const series = Array.from({length: maxItems}, (_, seriesIndex) => {
        const data = groups.map((group) => {
            const item = group.items[seriesIndex];
            if (!item) {
                return null;
            }

            const color = getDimensionColor(group.dimensionKey, seriesIndex);
            return {
                value: Number((item.accuracy * 100).toFixed(1)),
                accuracy: item.accuracy,
                dimensionLabel: group.dimensionLabel,
                valueLabel: item.label,
                itemStyle: {
                    color,
                    borderRadius: [8, 8, 0, 0],
                },
            } as AccuracyChartPoint;
        });

        return {
            name: `Item ${seriesIndex + 1}`,
            type: 'bar' as const,
            data,
            barMaxWidth: 64,
            itemStyle: {
                borderRadius: [8, 8, 0, 0],
            },
            emphasis: {focus: 'series' as const},
        };
    }).filter((serie) => serie.data.some((point) => point !== null));

    const hasSeriesData = series.some((serie) => serie.data.some((point) => point !== null));

    return {
        tooltip: {
            trigger: 'axis',
            axisPointer: {type: 'shadow'},
            formatter: (params: any) => {
                const points = Array.isArray(params) ? params.filter((item) => item?.data) : [params];
                if (!points.length) {
                    return '';
                }

                const dimension = points[0]?.data?.dimensionLabel ?? '';
                const lines = points.map((point: any) => {
                    const {valueLabel, accuracy} = point.data as AccuracyChartPoint;
                    const percent = typeof accuracy === 'number'
                        ? (accuracy * 100).toFixed(1)
                        : typeof point.value === 'number'
                            ? point.value.toFixed(1)
                            : point.value ?? '';
                    return `${valueLabel}: ${percent}%`;
                });

                const header = dimension ? `${dimension}<br/>` : '';
                return `${header}${lines.join('<br/>')}`;
            },
        },
        legend: {show: false},
        grid: {
            left: 48,
            right: 24,
            top: 10,
            bottom: 24,
        },
        xAxis: {
            type: 'category',
            data: xAxisCategories,
            axisLabel: {
                interval: 0,
                formatter: (value: string) => value.replace(/\s+/g, '\n'),
            },
        },
        yAxis: {
            type: 'value',
            min: 0,
            max: 100,
            axisLabel: {formatter: '{value}%'},
        },
        series: hasSeriesData ? series : [],
    };
};
