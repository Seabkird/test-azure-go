import { client } from '../../lib/api';

type HelloResponse = { message: string };

export const fetchHello = async (): Promise<HelloResponse> => {
  return await client<HelloResponse>('/api/user');
};