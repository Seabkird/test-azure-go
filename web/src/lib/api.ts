const apiUrl = import.meta.env.VITE_API_URL;

export const fetchHello = async () => {
  const res = await fetch(`${apiUrl}/api/user`); 
  if (!res.ok) throw new Error('Erreur réseau');
  return res.json();
};