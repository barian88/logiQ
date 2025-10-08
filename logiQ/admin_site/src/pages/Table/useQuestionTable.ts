import {useCallback, useEffect, useMemo, useState, type Key} from 'react';
import {message, type TablePaginationConfig, type TableProps} from 'antd';
import {keepPreviousData, useMutation, useQuery, useQueryClient} from '@tanstack/react-query';

import {deleteQuestions, fetchTableData, type ApiResponse} from './services.ts';
import type {Question, RawTableFilter, TableFilter} from './types.ts';
import {createJsonDownloadHandle} from '../../utils/download.ts';

const pickFirstFilterValue = (value: (string | number | boolean)[] | null | undefined): string =>
    Array.isArray(value) && value.length > 0 ? String(value[0]) : '';

const normalizeFilters = (rawFilters: RawTableFilter): TableFilter => ({
    category: pickFirstFilterValue(rawFilters.category),
    difficulty: pickFirstFilterValue(rawFilters.difficulty),
    type: pickFirstFilterValue(rawFilters.type),
});

export const useQuestionTable = () => {
    const [pagination, setPagination] = useState({current: 1, pageSize: 5});
    const [filters, setFilters] = useState<TableFilter>({category: '', difficulty: '', type: ''});
    const [selectedRowKeys, setSelectedRowKeys] = useState<Key[]>([]);

    const queryClient = useQueryClient();
    const [messageApi, contextHolder] = message.useMessage();

    const handleTableChange: TableProps<Question>['onChange'] = (
        _unusedPagination: TablePaginationConfig,
        rawFilters,
    ) => {
        const normalizedFilters = normalizeFilters(rawFilters as RawTableFilter);
        const filtersChanged =
            normalizedFilters.category !== filters.category ||
            normalizedFilters.difficulty !== filters.difficulty ||
            normalizedFilters.type !== filters.type;

        if (filtersChanged) {
            setFilters(normalizedFilters);
            setPagination((prev) => ({...prev, current: 1}));
        }
    };

    const handlePaginationChange = (page: number, pageSize: number) => {
        setPagination((prev) => {
            if (prev.current === page && prev.pageSize === pageSize) {
                return prev;
            }

            return {
                current: page,
                pageSize,
            };
        });
    };

    const handlePageSizeChange = (_current: number, pageSize: number) => {
        setPagination((prev) => {
            if (prev.pageSize === pageSize && prev.current === 1) {
                return prev;
            }

            return {
                current: 1,
                pageSize,
            };
        });
    };

    const {data, error, isFetching} = useQuery<ApiResponse, Error>({
        queryKey: ['tableData', pagination, filters],
        queryFn: () => fetchTableData(
            pagination.current,
            pagination.pageSize,
            filters.category || '',
            filters.difficulty || '',
            filters.type || ''
        ),
        placeholderData: keepPreviousData,
    });

    const {mutate: deleteSelectedRows, isPending: isDeleting} = useMutation({
        mutationFn: (ids: string[]) => deleteQuestions(ids),
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey: ['tableData']});
            setSelectedRowKeys([]);
            messageApi.success('successfully deleted');
        },
        onError: (err: unknown) => {
            const errorMessage = err instanceof Error ? err.message : 'failed to delete';
            messageApi.error(errorMessage);
        },
    });

    const tableData = useMemo<Question[]>(() => data?.list ?? [], [data]);
    const total = data?.total || 0;

    useEffect(() => {
        const validKeys = new Set(tableData.map((item) => String(item._id)));
        setSelectedRowKeys((prev) => prev.filter((key) => validKeys.has(String(key))));
    }, [tableData]);

    const rowSelection: TableProps<Question>['rowSelection'] = useMemo(
        () => ({
            selectedRowKeys,
            onChange: (keys) => setSelectedRowKeys(keys),
        }),
        [selectedRowKeys],
    );

    const handleDeleteConfirm = useCallback(() => {
        deleteSelectedRows(selectedRowKeys.map(String));
    }, [deleteSelectedRows, selectedRowKeys]);

    const handleExport = useCallback(() => {
        if (selectedRowKeys.length === 0) {
            return;
        }

        const selectedKeySet = new Set(selectedRowKeys.map(String));
        const selectedRows = tableData.filter((item) => selectedKeySet.has(String(item._id)));

        if (selectedRows.length === 0) {
            messageApi.warning('Selected rows are no longer available to export');
            setSelectedRowKeys([]);
            return;
        }

        const {url, revoke} = createJsonDownloadHandle(selectedRows);
        const link = document.createElement('a');
        link.href = url;
        const timestamp = new Date().toISOString().replace(/[:]/g, '-');
        link.download = `questions-${timestamp}.json`;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        revoke();
        messageApi.success(`Exported ${selectedRows.length} question(s)`);
    }, [messageApi, selectedRowKeys, tableData]);

    return {
        contextHolder,
        tableData,
        total,
        pagination,
        filters,
        rowSelection,
        isFetching,
        error,
        handleTableChange,
        handlePaginationChange,
        handlePageSizeChange,
        selectedRowKeys,
        isDeleting,
        handleDeleteConfirm,
        handleExport,
    };
};
