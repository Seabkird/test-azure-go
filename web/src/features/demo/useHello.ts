import { useQuery } from '@tanstack/react-query';
import { fetchHello } from './hello.api';

export const useHello = () => {
  return useQuery({
    queryKey: ['hello-message'], // La fameuse clé unique
    queryFn: fetchHello,         // La fonction à exécuter
    enabled: false,              // IMPORTANT : false = on attend un clic pour lancer
  });
};