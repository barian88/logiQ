import './table.css';

import {Button, Pagination, Popconfirm, Table} from 'antd';
import {DeleteOutlined, ExportOutlined, PlusCircleOutlined} from '@ant-design/icons';

import {questionTableColumns} from './columns.tsx';
import {useQuestionTable} from './useQuestionTable.ts';

function TablePage() {
    const {
        contextHolder,
        tableData,
        total,
        pagination,
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
    } = useQuestionTable();

    if (error) return <div>Error loading data: {error.message}</div>;

    return (
        <div className="table-container">
            {contextHolder}
            <Table
                dataSource={tableData}
                columns={questionTableColumns}
                rowKey="_id"
                loading={isFetching}
                scroll={{y: 'calc(100vh - 100px)'}}
                showSorterTooltip={{target: 'sorter-icon'}}
                rowSelection={rowSelection}
                pagination={false}
                onChange={handleTableChange}
            />
            <div className="table-footer-bar">
                <div className="table-footer-actions">
                    <Popconfirm
                        title="Are you sure to delete？"
                        okText="Yes"
                        cancelText="Cancel"
                        disabled={selectedRowKeys.length === 0}
                        okButtonProps={{loading: isDeleting}}
                        onConfirm={handleDeleteConfirm}
                    >
                        <Button
                            danger
                            type="primary"
                            icon={<DeleteOutlined />}
                            disabled={selectedRowKeys.length === 0}
                            loading={isDeleting}
                        >
                            Delete
                        </Button>
                    </Popconfirm>
                    <Button
                        type="primary"
                        icon={<ExportOutlined />}
                        disabled={selectedRowKeys.length === 0}
                        onClick={handleExport}
                    >
                        Export
                    </Button>
                    <Button type="primary" icon={<PlusCircleOutlined />}>
                        New
                    </Button>
                </div>
                <Pagination
                    current={pagination.current}
                    pageSize={pagination.pageSize}
                    total={total}
                    showSizeChanger
                    showQuickJumper
                    pageSizeOptions={[5, 10, 20, 50]}
                    onChange={handlePaginationChange}
                    onShowSizeChange={handlePageSizeChange}
                />
            </div>
        </div>
    );
}

export default TablePage;
