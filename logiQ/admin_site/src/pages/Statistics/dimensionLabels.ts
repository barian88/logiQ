import {CATEGORY_OPTIONS, DIFFICULTY_OPTIONS, TYPE_OPTIONS} from '../../constants/option';

export const DIMENSION_KEYS = ['category', 'difficulty', 'type'] as const;
export type DimensionKey = (typeof DIMENSION_KEYS)[number];

export const DIMENSION_LABELS: Record<DimensionKey, string> = {
    category: 'Category',
    difficulty: 'Difficulty',
    type: 'Type',
};

type Option = {label: string; value: string};

const buildLookup = (options: Option[]): Record<string, string> =>
    options.reduce<Record<string, string>>((map, option) => {
        map[option.value] = option.label;
        return map;
    }, {});

const OPTION_LOOKUP: Record<DimensionKey, Record<string, string>> = {
    category: buildLookup(CATEGORY_OPTIONS),
    difficulty: buildLookup(DIFFICULTY_OPTIONS),
    type: buildLookup(TYPE_OPTIONS),
};

const formatFallback = (value: string) =>
    value
        ?.replace(/[_-]/g, ' ')
        .replace(/\s+/g, ' ')
        .trim()
        .replace(/\b\w/g, (match) => match.toUpperCase()) || 'Unknown';

export const getDimensionValueLabel = (dimensionKey: DimensionKey, value: string): string => {
    return OPTION_LOOKUP[dimensionKey][value] ?? formatFallback(value);
};
