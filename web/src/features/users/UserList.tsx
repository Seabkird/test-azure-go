// src/features/users/UserList.tsx
import { Table, Container, Title, Paper, Badge } from '@mantine/core';

// Données fictives (Simule ta future API Go)
const MOCK_DATA = [
  { id: '1', name: 'Jean Dupont', email: 'jean@saas.com', role: 'ADMIN' },
  { id: '2', name: 'Alice Martin', email: 'alice@client.com', role: 'USER' },
  { id: '3', name: 'Bob Wilson', email: 'bob@tech.com', role: 'MANAGER' },
];

export function UserList() {
  const rows = MOCK_DATA.map((user) => (
    <Table.Tr key={user.id}>
      <Table.Td>{user.id}</Table.Td>
      <Table.Td fw={500}>{user.name}</Table.Td>
      <Table.Td>{user.email}</Table.Td>
      <Table.Td>
        {/* Exemple d'utilisation de composants Mantine pour le style */}
        <Badge color={user.role === 'ADMIN' ? 'blue' : 'gray'}>
          {user.role}
        </Badge>
      </Table.Td>
    </Table.Tr>
  ));

  return (
    <Container size="lg" py="xl">
      <Title order={2} mb="lg">Gestion des Utilisateurs</Title>
      <Paper withBorder shadow="sm" radius="md">
        <Table verticalSpacing="sm" striped>
          <Table.Thead>
            <Table.Tr>
              <Table.Th>ID</Table.Th>
              <Table.Th>Nom</Table.Th>
              <Table.Th>Email</Table.Th>
              <Table.Th>Rôle</Table.Th>
            </Table.Tr>
          </Table.Thead>
          <Table.Tbody>{rows}</Table.Tbody>
        </Table>
      </Paper>
    </Container>
  );
}