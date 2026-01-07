export const fetchHello = async () => {
  const res = await fetch('/api/user'); 
  if (!res.ok) throw new Error('Erreur réseau');
  return res.json();
};