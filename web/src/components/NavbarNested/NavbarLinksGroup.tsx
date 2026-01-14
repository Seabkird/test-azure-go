import { useState } from 'react';
import { Group, Box, Collapse, ThemeIcon, Text, UnstyledButton, rem } from '@mantine/core';
import { IconChevronRight } from '@tabler/icons-react';
import { Link, useLocation } from 'react-router-dom'; // IMPORTANT : Import Link
import classes from './NavbarLinksGroup.module.css';

interface LinksGroupProps {
  icon: React.FC<any>;
  label: string;
  initiallyOpened?: boolean;
  links?: { label: string; link: string }[];
  link?: string; // Ajout de la prop link pour les items simples
}

export function LinksGroup({ icon: Icon, label, initiallyOpened, links, link }: LinksGroupProps) {
  const hasLinks = Array.isArray(links);
  const [opened, setOpened] = useState(initiallyOpened || false);
  const location = useLocation(); // Pour gérer l'état actif (optionnel)

  // Cas 1 : C'est un groupe avec sous-menu (ex: Rapports)
  if (hasLinks) {
      const items = (links || []).map((link) => (
        <Text
            component={Link}
            to={link.link}
            className={classes.link}
            key={link.label}
            data-active={location.pathname === link.link || undefined} 
        >
            {link.label}
        </Text>
        ));

      return (
        <>
          <UnstyledButton onClick={() => setOpened((o) => !o)} className={classes.control}>
            <Group justify="space-between" gap={0}>
              <Box style={{ display: 'flex', alignItems: 'center' }}>
                <ThemeIcon variant="light" size={30}>
                  <Icon style={{ width: rem(18), height: rem(18) }} />
                </ThemeIcon>
                <Box ml="md">{label}</Box>
              </Box>
              {hasLinks && (
                <IconChevronRight
                  className={classes.chevron}
                  stroke={1.5}
                  style={{
                    width: rem(16),
                    height: rem(16),
                    transform: opened ? 'rotate(-90deg)' : 'none',
                  }}
                />
              )}
            </Group>
          </UnstyledButton>
          <Collapse in={opened}>{items}</Collapse>
        </>
      );
  }

  // Cas 2 : C'est un lien simple (ex: Accueil, Utilisateurs)
  return (
    <UnstyledButton
      component={Link} // IMPORTANT : Utiliser Link ici aussi
      to={link || '#'} // IMPORTANT : vers le lien défini dans mockdata
      className={classes.control}
      data-active={location.pathname === link || undefined}
    >
      <Group justify="space-between" gap={0}>
        <Box style={{ display: 'flex', alignItems: 'center' }}>
          <ThemeIcon variant="light" size={30}>
            <Icon style={{ width: rem(18), height: rem(18) }} />
          </ThemeIcon>
          <Box ml="md">{label}</Box>
        </Box>
      </Group>
    </UnstyledButton>
  );
}