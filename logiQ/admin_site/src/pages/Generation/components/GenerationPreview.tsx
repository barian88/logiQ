import React, {useMemo} from 'react';
import {Alert, Button, Card, Divider, List, Space, Typography} from 'antd';
import {DownloadOutlined} from '@ant-design/icons';
import ReactECharts from 'echarts-for-react';

import type {Question} from '../../Table/types.ts';

const {Paragraph, Text} = Typography;

export interface DistributionItem {
    value: string;
    label: string;
    count: number;
}

export interface GenerationSummary {
    total: number;
    category: DistributionItem[];
    difficulty: DistributionItem[];
    type: DistributionItem[];
}

interface GenerationPreviewProps {
    summary: GenerationSummary;
    previewItems: Question[];
    generatedAt: Date | null;
    downloadUrl?: string;
}

const buildSeriesData = (items: DistributionItem[]) =>
    items
        .filter((item) => item.count > 0)
        .map((item) => {
            const abbreviation = item.label
                .split(/\s|\//)
                .filter(Boolean)
                .map((token) => token[0]?.toUpperCase() ?? '')
                .join('');

            return {
                value: item.count,
                name: abbreviation || item.label,
                originalLabel: item.label,
            };
        });

const buildNestedChartOption = (summary: GenerationSummary) => {
    const difficultyData = buildSeriesData(summary.difficulty);
    const typeData = buildSeriesData(summary.type);
    const categoryData = buildSeriesData(summary.category);

    const legendLabels = new Map<string, string>();
    [...difficultyData, ...typeData, ...categoryData].forEach((item) => {
        legendLabels.set(item.name, item.originalLabel);
    });

    const tooltipFormatter = (params: any) => {
        const {seriesName, data} = params;
        const label = data?.originalLabel ?? params.name;
        return `${seriesName}<br/>${label}: ${params.value} (${params.percent}%)`;
    };

    const legendFormatter = (name: string) => legendLabels.get(name) ?? name;

    const series = [
        {
            name: 'Difficulty',
            type: 'pie',
            center: ['40%', '50%'],
            radius: [0, '30%'],
            label: {
                position: 'inner',
                formatter: '{b}\n{c}',
            },
            data: difficultyData,
        },
        {
            name: 'Type',
            type: 'pie',
            center: ['40%', '50%'],
            radius: ['40%', '60%'],
            label: {
                position: 'inner',
                formatter: '{b}\n{c}',
            },
            data: typeData,
        },
        {
            name: 'Category',
            type: 'pie',
            center: ['40%', '50%'],
            radius: ['70%', '90%'],
            label: {
                position: 'inner',
                formatter: '{b}\n{c}',
            },
            data: categoryData,
        },
    ].filter((layer) => Array.isArray(layer.data) && layer.data.length > 0);

    return {
        tooltip: {
            trigger: 'item',
            formatter: tooltipFormatter,
        },
        legend: {
            right: 0,
            top: 'center',
            orient: 'vertical',
            itemWidth: 10,
            padding: [20, 0, 0, 0],
            type: 'scroll',
            formatter: legendFormatter,
        },
        series,
    };
};

const GenerationPreview: React.FC<GenerationPreviewProps> = ({
    summary,
    previewItems,
    generatedAt,
    downloadUrl,
}) => {
    if (summary.total === 0) {
        return (
            <Alert
                type="info"
                showIcon
                message="Preview"
                description="Submit the form to see the latest generation summary and download the results."
            />
        );
    }

    const chartOption = useMemo(() => buildNestedChartOption(summary), [summary]);

    return (
        <Card
            title="Generation Preview"
            extra={
                <Button
                    type="link"
                    icon={<DownloadOutlined />}
                    href={downloadUrl}
                    download="generated-questions.json"
                    disabled={!downloadUrl}
                >
                    Download JSON
                </Button>
            }
        >
            <Space direction="vertical" size={5} style={{width: '100%', flex: 1, overflow: 'hidden'}}>
                {generatedAt && (
                    <Text type="secondary">Generated at: {generatedAt.toLocaleString()}</Text>
                )}
                <div>
                    <Text type={"secondary"}> Total Questions: </Text>
                    <Text strong>{summary.total}</Text>
                    </div>
                <div style={{flex: 1, minHeight: 0, display: 'flex', flexDirection: 'column'}}>
                    <div>
                        {chartOption.series.length > 0 ? (
                            <ReactECharts
                                option={chartOption}
                                style={{height: '260px', margin: '0 auto'}}
                                opts={{renderer: 'svg'}}
                            />
                        ) : (
                            <Paragraph type="secondary" style={{marginBottom: 0}}>
                                Not enough data to draw the chart yet.
                            </Paragraph>
                        )}
                    </div>
                    {previewItems.length > 0 && (
                        <div style={{marginTop: 12, minHeight: 0, display: 'flex', flexDirection: 'column'}}>
                            <Divider style={{margin: '12px 0'}} />
                            <Text strong>Sample Questions</Text>
                            <List
                                rowKey={(item) => item._id}
                                size="small"
                                dataSource={previewItems}
                                style={{maxHeight: 200, overflowY: 'auto'}}
                                renderItem={(item, index) => (
                                    <List.Item>
                                        <Space direction="horizontal" size={8} align="start">
                                            <Text strong>{index + 1}</Text>
                                            <Text>{item.question_text}</Text>
                                        </Space>
                                    </List.Item>
                                )}
                            />
                        </div>
                    )}
                </div>
            </Space>
        </Card>
    );
};

export default GenerationPreview;
