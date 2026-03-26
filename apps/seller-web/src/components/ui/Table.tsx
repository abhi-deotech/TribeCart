
import React from 'react';

interface Column<T> {
  header: string;
  accessor: keyof T | ((row: T) => React.ReactNode);
  className?: string;
}

interface TableProps<T> {
  columns: Column<T>[];
  data: T[];
  onRowClick?: (row: T) => void;
  keyExtractor?: (row: T, index: number) => string;
  className?: string;
  getRowProps?: (row: T) => React.HTMLAttributes<HTMLTableRowElement>;
  isLoading?: boolean;
  emptyStateMessage?: string;
}

const Table = <T,>({
  columns = [],
  data = [],
  onRowClick,
  keyExtractor = (row, i) => (row as any).id || (row as any)._id || `row-${i}`,
  className = '',
  getRowProps = () => ({}),
  isLoading = false,
  emptyStateMessage = 'No data available',
}: TableProps<T>) => {
  return (
    <div className={`overflow-x-auto rounded-xl ${className}`}>
      <table className="min-w-full divide-y divide-gray-200">
        <thead className="bg-gray-50">
          <tr>
            {columns.map((col, i) => (
              <th
                key={col.header || `col-${i}`}
                className={`px-5 py-3 text-left text-sm font-semibold text-gray-700 ${col.className || ''}`}
              >
                {col.header}
              </th>
            ))}
          </tr>
        </thead>

        <tbody className="bg-white divide-y divide-gray-200">
          {isLoading ? (
            <tr>
              <td colSpan={columns.length} className="px-6 py-4 text-center text-sm text-gray-500">
                Loading...
              </td>
            </tr>
          ) : data.length === 0 ? (
            <tr>
              <td
                colSpan={columns.length}
                className="px-6 py-4 text-center text-sm text-gray-500"
              >
                {emptyStateMessage}
              </td>
            </tr>
          ) : (
            data.map((row, rowIdx) => {
              const { className: rowClassName, ...restRowProps } = getRowProps(row);
              const trClassName = [
                onRowClick ? 'cursor-pointer hover:bg-gray-50' : '',
                rowClassName,
              ]
                .filter(Boolean)
                .join(' ');
              return (
                <tr
                  key={keyExtractor(row, rowIdx)}
                  className={trClassName}
                  onClick={() => onRowClick?.(row)}
                  {...restRowProps}
                >
                  {columns.map((col, colIdx) => {
                    const content =
                      typeof col.accessor === 'function' ? col.accessor(row) : row[col.accessor];

                    return (
                      <td
                        key={`${keyExtractor(row, rowIdx)}-${String(col.accessor) || colIdx}`}
                        className={`px-6 py-4 text-sm text-gray-800 ${col.className || ''} relative`}
                      >
                        {content}
                      </td>
                    );
                  })}
                </tr>
              );
            })
          )}
        </tbody>
      </table>
    </div>
  );
};

export default Table;
