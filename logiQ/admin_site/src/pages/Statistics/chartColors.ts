import type {DimensionKey} from './dimensionLabels';

export const STATISTICS_COLOR_PALETTE = [
    '#5470C6', // blue
    '#91CC75', // green
    '#EE6666', // red
    '#FAC858', // yellow
    '#73C0DE', // cyan
    '#3BA272', // teal
    '#FC8452', // orange
    '#9A60B4', // purple
    '#EA7CCC', // pink
];

const DIMENSION_COLOR_OFFSETS: Record<DimensionKey, number> = {
    category: 0,
    difficulty: 3,
    type: 6,
};

export const getDimensionColor = (dimensionKey: DimensionKey, itemIndex: number): string => {
    const offset = DIMENSION_COLOR_OFFSETS[dimensionKey] ?? 0;
    const paletteIndex = (offset + itemIndex) % STATISTICS_COLOR_PALETTE.length;
    return STATISTICS_COLOR_PALETTE[paletteIndex];
};

export const getDimensionColorPalette = (dimensionKey: DimensionKey, itemCount: number): string[] => {
    return Array.from({length: itemCount}, (_, index) => getDimensionColor(dimensionKey, index));
};

