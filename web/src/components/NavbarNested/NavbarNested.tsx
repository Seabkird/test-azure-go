import { Group, Code, ScrollArea } from '@mantine/core';
import {
  IconNotes,
  IconCalendarStats,
  IconGauge,
  IconPresentationAnalytics,
  IconFileAnalytics,
  IconAdjustments,
  IconLock,
  IconUsers, // J'ai ajouté l'icone Users
} from '@tabler/icons-react';
import { UserButton } from './UserButton'; // Assure-toi que le chemin est bon
import { LinksGroup } from './NavbarLinksGroup'; // Assure-toi que le chemin est bon
import { Logo } from './Logo'; // Assure-toi que le chemin est bon
import classes from './NavbarNested.module.css';

// C'est ICI que tu définis tes routes
const mockdata = [
  { label: 'Accueil', icon: IconGauge, link: '/' }, // Lien vers Home
  { label: 'Utilisateurs', icon: IconUsers, link: '/users' }, // Lien vers Users
  
  // Exemple de menu déroulant (gardé pour l'exemple)
  {
    label: 'Rapports',
    icon: IconNotes,
    initiallyOpened: false,
    links: [
      { label: 'Vue d\'ensemble', link: '/reports/overview' },
      { label: 'Prévisions', link: '/reports/forecasts' },
    ],
  },
  { label: 'Analytics', icon: IconPresentationAnalytics },
  { label: 'Settings', icon: IconAdjustments },
];

export function NavbarNested() {
  const links = mockdata.map((item) => <LinksGroup {...item} key={item.label} />);

  return (
    <nav className={classes.navbar}>
      <div className={classes.header}>
        <Group justify="space-between">
          <Logo style={{ width: 120 }} />
          <Code fw={700}>v1.0.0</Code>
        </Group>
      </div>

      <ScrollArea className={classes.links}>
        <div className={classes.linksInner}>{links}</div>
      </ScrollArea>

      <div className={classes.footer}>
        <UserButton />
      </div>
    </nav>
  );
}