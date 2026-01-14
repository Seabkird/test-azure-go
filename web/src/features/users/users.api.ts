import { client } from '../../lib/api';
import { UserApiResponse } from './users.types';

// La fonction métier spécifique
export const fetchUsers = async (): Promise<UserApiResponse[]> => {
  return client<UserApiResponse[]>('/api/users');
};