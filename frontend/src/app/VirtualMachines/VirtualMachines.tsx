import * as React from 'react';
import VirtualMachineList from '../components/VirtualMachineList';
import BasicHint from '../components/basichint';
import {
  Breadcrumb, 
  BreadcrumbItem,
  PageSectionVariants,
  Page,
  TextContent,
  PageSection,
  Text
} from '@patternfly/react-core';
import { TableComposable, Thead, Tbody, Tr, Th, Td, Caption } from '@patternfly/react-table';

const columns = ['Name', 'Latest Event ID', 'Role', 'Workspaces', 'Last commit'];
const rows = [
  ['one', 'two', 'three', 'four', 'five'],
  ['one - 2', null, null, 'four - 2', 'five - 2'],
  ['one - 3', 'two - 3', 'three - 3', 'four - 3', 'five - 3']
];


// eslint-disable-next-line prefer-const
const VirtualMachines: React.FunctionComponent = () => (
<Page>
  <PageSection variant={PageSectionVariants.light}>
    <TextContent>
      <Text component="h1">Virtual Machines</Text>
    </TextContent>
  </PageSection>
  <PageSection>
      <TableComposable
        aria-label="Simple table"
        variant="compact"
      >
        <Thead>
          <Tr>
            {columns.map((column, columnIndex) => (
              <Th key={columnIndex}>{column}</Th>
            ))}
          </Tr>
        </Thead>
        <Tbody>
          {rows.map((row, rowIndex) => (
            <Tr key={rowIndex}>
              {row.map((cell, cellIndex) => (
                <Td key={`${rowIndex}_${cellIndex}`} dataLabel={columns[cellIndex]}>
                  {cell}
                </Td>
              ))}
            </Tr>
          ))}
        </Tbody>
      </TableComposable>
  </PageSection>
</Page>
)

export { VirtualMachines };
