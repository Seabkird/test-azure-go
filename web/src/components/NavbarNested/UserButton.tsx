import { UnstyledButton, Group, Avatar, Text, rem } from '@mantine/core';
import { IconChevronRight } from '@tabler/icons-react';
import classes from './UserButton.module.css';

interface UserButtonProps extends React.ComponentPropsWithoutRef<'button'> {
  image?: string;
  name?: string;
  email?: string;
}

export function UserButton({ image, name, email, ...others }: UserButtonProps) {
  // Données par défaut (Mock) - Analogie Backend: comme des "fixtures"
  const defaultImage = 'https://raw.githubusercontent.com/mantinedev/mantine/master/.demo/avatars/avatar-8.png';
  const defaultName = 'Harriette Spoonlicker';
  const defaultEmail = 'hspoonlicker@outlook.com';

  return (
    <UnstyledButton className={classes.user} {...others}>
      <Group>
        <Avatar src={image || defaultImage} radius="xl" />

        <div style={{ flex: 1 }}>
          <Text size="sm" fw={500}>
            {name || defaultName}
          </Text>

          <Text c="dimmed" size="xs">
            {email || defaultEmail}
          </Text>
        </div>

        <IconChevronRight style={{ width: rem(14), height: rem(14) }} stroke={1.5} />
      </Group>
    </UnstyledButton>
  );
}