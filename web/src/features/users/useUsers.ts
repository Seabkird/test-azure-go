import { useQuery } from '@tanstack/react-query';
import { fetchUsers } from './users.api';
import { UserApiResponse } from './users.types';

export const useUsers = () => {
  return useQuery<UserApiResponse[]>({
    queryKey: ['users', 'list'], // Clé plus descriptive
    queryFn: fetchUsers,
  });
};