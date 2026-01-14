// src/features/users/UserList.tsx
import { Table, Container, Title, Paper, Badge, Loader, Center, Alert, Text } from '@mantine/core';
import { IconAlertCircle } from '@tabler/icons-react';
import { useUsers } from './useUsers';
import { UserApiResponse } from './users.types';

export function UserList() {
  // 2. APPEL DU HOOK TANSTACK QUERY
  // useQuery retourne un objet riche. On déstructure les propriétés dont on a besoin.
  // On renomme 'data' en 'users' pour plus de clarté dans ce composant.
  const {
    data: users,     // Sera de type UserApiResponse[] | undefined
    isLoading,       // true tant que la requête HTTP n'est pas finie
    isError,         // true si le backend renvoie une 4xx/5xx ou si le réseau échoue
    error            // L'objet erreur complet
  } = useUsers();

  // 3. GESTION DES ÉTATS (Pattern indispensable en React)

  // État de chargement : Afficher un spinner Mantine
  if (isLoading) {
    return (
      <Container size="lg" py="xl">
        <Center h={400}>
           {/* Vous pouvez personnaliser le loader */}
           <Loader size="xl" variant="dots" color="blue" />
        </Center>
      </Container>
    );
  }

  // État d'erreur : Afficher une alerte explicite
  if (isError) {
    // 'error' est ici un objet Error (ou celui défini par votre fetcher)
    // Pour un dev backend, c'est l'équivalent de catcher une exception et d'afficher un message à l'utilisateur.
    return (
      <Container size="lg" py="xl">
         <Alert icon={<IconAlertCircle size={16} />} title="Erreur de chargement" color="red" variant="filled">
           Impossible de récupérer la liste des utilisateurs. Veuillez réessayer plus tard.
           {/* En dev, vous pouvez afficher le message technique : */}
           {/* <Text size="sm" mt="sm">{error instanceof Error ? error.message : 'Erreur inconnue'}</Text> */}
         </Alert>
      </Container>
    );
  }

  // 4. SÉCURISATION DES DONNÉES
  // Si isLoading est false et isError est false, 'users' devrait être défini.
  // Cependant, TypeScript sait qu'il peut encore être 'undefined' avant le tout premier montage réussi.
  // On utilise l'opérateur "nullish coalescing" (??) ou un "OR" logique (||) pour garantir un tableau.
  const safeUsersList: UserApiResponse[] = users || [];


  // Gestion du cas "Pas de données" (tableau vide renvoyé par l'API)
  if (safeUsersList.length === 0) {
      return (
        <Container size="lg" py="xl">
           <Title order={2} mb="lg">Gestion des Utilisateurs</Title>
            <Paper withBorder p="xl" ta="center" bg="gray.0">
               <Text c="dimmed" fs="italic">Aucun utilisateur trouvé dans la base de données.</Text>
            </Paper>
        </Container>
      );
  }

  // 5. MAPPING DES DONNÉES VERS L'UI (Rendu du tableau)
  // TypeScript va ici vérifier que user.id, user.name, etc., existent bien dans UserApiResponse
  const rows = safeUsersList.map((user) => (
    <Table.Tr key={user.id}>
      <Table.Td fw={500}>#{user.id.substring(0, 8)}</Table.Td> {/* Exemple: raccourcir un UUID */}
      <Table.Td fw={700}>{user.nom}</Table.Td>
      <Table.Td fw={700}>{user.prenom}</Table.Td>
      <Table.Td>{user.email}</Table.Td>
      <Table.Td>
        {/* Adaptez la logique des couleurs selon les valeurs réelles de votre enum 'role' */}
        <Badge
          color={
            user.role === 'ADMIN' ? 'red' :
            user.role === 'MANAGER' ? 'blue' : 'green'
          }
          variant="light"
        >
          {user.role}
        </Badge>
      </Table.Td>
    </Table.Tr>
  ));


  // Rendu final (Structure principale)
  return (
    <Container size="lg" py="xl">
      <Title order={2} mb="lg">Gestion des Utilisateurs</Title>
      <Paper withBorder shadow="md" radius="md" p="md">
        <Table verticalSpacing="sm" highlightOnHover>
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