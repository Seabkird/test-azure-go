/**
 * ARCHITECTURE NOTE: LAYOUT PATTERN
 * ---------------------------------
 * Ce composant agit comme le "Master Page" ou le "Base Template" en Backend.
 * 
 * Rôles :
 * 1. Structure Globale : Définit le squelette immuable (Header, Navbar, Footer).
 * 2. Persistance : Ces éléments ne sont pas rechargés lors de la navigation.
 * 3. Injection de Contenu (Slot) : Le composant <Outlet /> agit comme un placeholder dynamique.
 *    C'est là que le Routeur injectera la page demandée (ex: Home, Users).
 * 
 * Usage :
 * Il est utilisé comme "Parent Route" dans Router.tsx.
 */

import { AppShell, Burger, Group, NavLink } from '@mantine/core'; // J'ai ajouté NavLink
import { useDisclosure } from '@mantine/hooks';
import { Outlet, useLocation, Link } from 'react-router-dom'; // J'ai ajouté Link et useLocation
import { NavbarNested } from '../NavbarNested/NavbarNested';

export function MainLayout() {
  const [opened, { toggle }] = useDisclosure();
  const location = useLocation(); // Pour savoir sur quelle page on est (active state)

  return (
    <AppShell
      header={{ height: 60 }}
      navbar={{
        width: 300,
        breakpoint: 'sm',
        collapsed: { mobile: !opened },
      }}
      padding="md"
    >
      <AppShell.Header>
        <Group h="100%" px="md">
          <Burger opened={opened} onClick={toggle} hiddenFrom="sm" size="sm" />
          <div style={{ fontWeight: 'bold' }}>Mon SaaS</div>
        </Group>
      </AppShell.Header>

      {/* On supprime le padding ici car NavbarNested gère son propre layout */}
      <AppShell.Navbar p={0}> 
        <NavbarNested />
      </AppShell.Navbar>

      <AppShell.Main>
        <Outlet /> 
      </AppShell.Main>
    </AppShell>
  );
}