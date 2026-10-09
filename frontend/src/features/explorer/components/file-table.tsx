import type { FileItem } from "#/client/gen/index.ts";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "#/components/ui/table.tsx";
import { useTable } from "@tanstack/react-table";
import { columns } from "./columns";
import { features } from "./features";

type FileTableProps = {
  data: FileItem[];
};

export const FileTable = ({ data }: FileTableProps) => {
  const table = useTable({
    key: "file-table",
    features,
    columns,
    data,
  });

  return (
    <div className="bg-card shadow border-2">
      <Table className="w-full text-base">
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                const hasExplicitSize = header.column.columnDef.size !== undefined;
                return (
                  <TableHead
                    className="font-bold"
                    key={header.id}
                    style={{
                      width: hasExplicitSize ? header.column.getSize() : "auto",
                    }}
                  >
                    {header.isPlaceholder ? null : <table.FlexRender header={header} />}
                  </TableHead>
                );
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows?.length ? (
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id} data-state={row.getIsSelected() && "selected"}>
                {row.getAllCells().map((cell) => {
                  const hasExplicitSize = cell.column.columnDef.size !== undefined;
                  return (
                    <TableCell
                      key={cell.id}
                      style={{
                        width: hasExplicitSize ? cell.column.getSize() : "auto",
                      }}
                    >
                      <table.FlexRender cell={cell} />
                    </TableCell>
                  );
                })}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                No results.
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
};
