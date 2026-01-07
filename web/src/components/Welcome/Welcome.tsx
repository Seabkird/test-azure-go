import { Title, Text, Anchor, Button, Group, Loader  } from '@mantine/core';
import { Link } from 'react-router-dom'; // <--- Import important
import { useHello } from '../../features/demo/useHello';
import classes from './Welcome.module.css';

export function Welcome() {
  const { data, isLoading, error, refetch } = useHello();

  return (
    <>
      <Title className={classes.title} ta="center" mt={100}>
        Welcome to{' '}
        <Text inherit variant="gradient" component="span" gradient={{ from: 'pink', to: 'yellow' }}>
          Mantine
        </Text>
      </Title>
      
      <Group justify="center" mt="xl">
        <Button 
          component={Link} // Transforme le bouton Mantine en lien Router
          to="/users"      // L'URL définie dans Router.tsx
          size="lg"
          variant="filled"
        >
          Voir les utilisateurs
        </Button>
      </Group>
        <Group justify="center" mt="xl">
        <Button onClick={() => refetch()}>
          Dire Bonjour au Backend
        </Button>

        {isLoading && <Loader />}
        
        {error && <Text c="red">Erreur : {error.message}</Text>}

        {/* Si data existe, on l'affiche */}
        {data && (
          <Text mt="xl" fw={700} c="blue">
            Réponse du serveur : {data.message}
          </Text>
        )}
      </Group>

      {/* ------------------------- */}
    </>
  );
}