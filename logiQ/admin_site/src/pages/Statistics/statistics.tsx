import {useMemo} from 'react';
import type {CSSProperties} from 'react';
import type {EChartsOption} from 'echarts';
import ReactECharts from 'echarts-for-react';
import {Alert, Card, Col, Row} from 'antd';
import {useQuery} from '@tanstack/react-query';

import './statistics.css';
import {buildAccuracyOption, toAccuracyChartGroups} from './accuracyChart';
import {buildPieOption} from './distributionChart';
import {
    type AccuracyDistributions,
    type DimensionDistributions,
    getAccuracyDistributions,
    getDimensionDistributions,
} from './services.ts';

const chartStyle: CSSProperties = {
    height: 230,
};

const pieChartStyle: CSSProperties = {
    height: 220,
};

const emptyStateStyle: CSSProperties = {
    height: 220,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    color: '#999',
    background: '#f5f5f5',
    borderRadius: 8,
};

const StatisticsPage = () => {
    const {
        data: accuracyDistributions,
        isLoading: isLoadingAccuracy,
        isError: isAccuracyError,
        error: accuracyError,
    } = useQuery<AccuracyDistributions>({
        queryKey: ['dimensionAccuracy'],
        queryFn: getAccuracyDistributions,
    });

    const {
        data: dimensionDistributions,
        isLoading: isLoadingDistributions,
        isError: isDistributionError,
        error: distributionError,
    } = useQuery<DimensionDistributions>({
        queryKey: ['dimensionDistributions'],
        queryFn: getDimensionDistributions,
    });

    const groupedAccuracy = useMemo(
        () => toAccuracyChartGroups(accuracyDistributions),
        [accuracyDistributions],
    );

    const accuracyOption = useMemo(() => {
        return groupedAccuracy.length ? buildAccuracyOption(groupedAccuracy) : undefined;
    }, [groupedAccuracy]);

    const categoryData = dimensionDistributions?.category;
    const difficultyData = dimensionDistributions?.difficulty;
    const typeData = dimensionDistributions?.type;

    const categoryPieOption = useMemo(
        () => (categoryData && categoryData.length ? buildPieOption('category', categoryData) : undefined),
        [categoryData],
    );

    const difficultyPieOption = useMemo(
        () => (difficultyData && difficultyData.length ? buildPieOption('difficulty', difficultyData) : undefined),
        [difficultyData],
    );

    const typePieOption = useMemo(
        () => (typeData && typeData.length ? buildPieOption('type', typeData) : undefined),
        [typeData],
    );

    const renderPieSection = (title: string, option?: EChartsOption, hasData = false) => (
        <div className="statistics-pie-item" key={title}>
            {hasData && option ? (
                <ReactECharts option={option} style={pieChartStyle} opts={{renderer: 'svg'}} />
            ) : (
                <div style={emptyStateStyle}>No data yet</div>
            )}
        </div>
    );

    return (
        <div className={'container'}>
            <Row gutter={[24, 24]}>
                <Col span={24}>
                    <Card title="Question Accuracy" bodyStyle={{padding: 24}} loading={isLoadingAccuracy}>
                        {isAccuracyError ? (
                            <Alert
                                type="error"
                                message="Unable to load accuracy data"
                                description={accuracyError instanceof Error ? accuracyError.message : 'Unknown error'}
                                showIcon
                            />
                        ) : accuracyOption ? (
                            <ReactECharts option={accuracyOption} style={chartStyle} opts={{renderer: 'svg'}} />
                        ) : (
                            <div style={emptyStateStyle}>No accuracy data yet</div>
                        )}
                    </Card>
                </Col>
            </Row>
            <Row gutter={[24, 24]}>
                <Col span={24}>
                    <Card
                        title="Question Distribution"
                        bodyStyle={{padding: 24}}
                        loading={isLoadingDistributions}
                    >
                        {isDistributionError ? (
                            <Alert
                                type="error"
                                message="Unable to load dimension distribution"
                                description={distributionError instanceof Error ? distributionError.message : 'Unknown error'}
                                showIcon
                            />
                        ) : (
                            <div className="statistics-pie-grid">
                                {renderPieSection('By Category', categoryPieOption, Boolean(categoryData?.length))}
                                {renderPieSection('By Difficulty', difficultyPieOption, Boolean(difficultyData?.length))}
                                {renderPieSection('By Type', typePieOption, Boolean(typeData?.length))}
                            </div>
                        )}
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default StatisticsPage;
