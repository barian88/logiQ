import type {EChartsOption} from 'echarts';

import type {DimensionPortion} from './services';
import type {DimensionKey} from './dimensionLabels';
import {getDimensionValueLabel} from './dimensionLabels';
import {getDimensionColorPalette} from './chartColors';

export const buildPieOption = (dimensionKey: DimensionKey, data: DimensionPortion[]): EChartsOption => {
    const seriesData = data.map((item) => ({
        name: getDimensionValueLabel(dimensionKey, item.value),
        value: item.count,
        portion: item.portion,
    }));

    const colorPalette = getDimensionColorPalette(dimensionKey, seriesData.length);

    return {
        color: colorPalette,
        tooltip: {
            trigger: 'item',
            formatter: (params: any) => {
                const portion = typeof params?.data?.portion === 'number'
                    ? `${Math.round(params.data.portion * 100)}%`
                    : `${params.percent}%`;
                return `${params.name}: ${params.value} (${portion})`;
            },
        },
        series: [
            {
                name: 'Distribution',
                type: 'pie',
                radius: ['0', '80%'],
                center: ['50%', '45%'],
                label: {
                    formatter: (params: any) => `${params.name}\n${params.percent}%`,
                },
                data: seriesData,
            },
        ],
    };
};
